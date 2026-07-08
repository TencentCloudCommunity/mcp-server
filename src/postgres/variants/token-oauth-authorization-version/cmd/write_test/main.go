package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"postgres_server/security"
)

type StepRunner struct {
	Name    string
	Tool    string
	Prepare func(*Runner, StepConfig) (map[string]interface{}, string, error)
	Verify  func(*Runner, StepConfig, map[string]interface{}, string) error
}

type Runner struct {
	plan       *TestPlan
	client     *client.Client
	ctx        context.Context
	targetMeta *InstanceMeta
	passed     int
	failed     int
	skipped    int
	timeout    time.Duration
}

type InstanceMeta struct {
	ID             string
	Name           string
	Status         string
	Zone           string
	VpcId          string
	SubnetId       string
	SpecCode       string
	DBMajorVersion string
}

func main() {
	configPath := flag.String("config", "scripts/full_test_plan.yaml", "写操作批量验证配置文件路径")
	onlyStep := flag.String("only", "", "仅执行指定 step 名称")
	flag.Parse()

	plan, err := LoadPlan(*configPath)
	if err != nil {
		fmt.Println("load config error:", err)
		os.Exit(1)
	}

	runner, err := NewRunner(plan)
	if err != nil {
		fmt.Println("init runner error:", err)
		os.Exit(1)
	}
	defer runner.Close()

	fmt.Println("========== CONFIG DRIVEN WRITE TEST ==========")
	fmt.Printf("Config: %s\n", *configPath)
	fmt.Printf("Target: %s (%s)\n", plan.Target.DBInstanceID, plan.Target.Region)
	fmt.Printf("Server: %s\n", plan.Server.URL)
	fmt.Println("说明：只有 enabled=true 且 approved=true 的步骤才会真实执行")
	fmt.Println()

	if err := runner.Run(orderedSteps(), *onlyStep); err != nil {
		fmt.Println("run error:", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("========== SUMMARY ==========")
	fmt.Printf("Passed: %d, Failed: %d, Skipped: %d\n", runner.passed, runner.failed, runner.skipped)
	if runner.failed > 0 {
		os.Exit(1)
	}
}

func NewRunner(plan *TestPlan) (*Runner, error) {
	c, err := client.NewSSEMCPClient(plan.Server.URL, security.MCPClientOptionsFromEnv()...)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	if err := c.Start(ctx); err != nil {
		return nil, err
	}
	if _, err := c.Initialize(ctx, mcp.InitializeRequest{}); err != nil {
		_ = c.Close()
		return nil, err
	}
	return &Runner{
		plan:    plan,
		client:  c,
		ctx:     ctx,
		timeout: time.Duration(plan.Options.TimeoutSeconds) * time.Second,
	}, nil
}

func (r *Runner) Close() {
	if r.client != nil {
		r.client.Close()
	}
}

func orderedSteps() []StepRunner {
	return []StepRunner{
		{Name: "ModifyDBInstanceName", Tool: "ModifyDBInstanceName", Prepare: prepareInstanceStep, Verify: verifyInstanceName},
		{Name: "CreateAccount", Tool: "CreateAccount", Prepare: prepareInstanceStep, Verify: verifyAccountCreated},
		{Name: "ResetAccountPassword", Tool: "ResetAccountPassword", Prepare: prepareInstanceStep},
		{Name: "ModifyAccountPrivilegesGrant", Tool: "ModifyAccountPrivileges", Prepare: prepareInstanceStep, Verify: verifyPrivilegeGrant},
		{Name: "ModifyAccountPrivilegesRevoke", Tool: "ModifyAccountPrivileges", Prepare: prepareInstanceStep, Verify: verifyPrivilegeRevoke},
		{Name: "DeleteAccount", Tool: "DeleteAccount", Prepare: prepareInstanceStep, Verify: verifyAccountDeleted},
		{Name: "CreateDatabase", Tool: "CreateDatabase", Prepare: prepareInstanceStep, Verify: verifyDatabaseCreated},
		{Name: "ModifyDatabaseOwner", Tool: "ModifyDatabaseOwner", Prepare: prepareInstanceStep, Verify: verifyDatabaseOwner},
		{Name: "DescribeBackupDownloadURL", Tool: "DescribeBackupDownloadURL", Prepare: prepareBackupDownloadURL},
		{Name: "CreateBaseBackup", Tool: "CreateBaseBackup", Prepare: prepareInstanceStep},
		{Name: "OpenDBExtranetAccess", Tool: "OpenDBExtranetAccess", Prepare: prepareInstanceStep},
		{Name: "CloseDBExtranetAccess", Tool: "CloseDBExtranetAccess", Prepare: prepareInstanceStep},
		{Name: "ModifyDBInstanceSecurityGroups", Tool: "ModifyDBInstanceSecurityGroups", Prepare: prepareModifySecurityGroups, Verify: verifySecurityGroups},
		{Name: "UpgradeDBInstanceKernelVersion", Tool: "UpgradeDBInstanceKernelVersion", Prepare: prepareInstanceStep},
		{Name: "RestartDBInstance", Tool: "RestartDBInstance", Prepare: prepareInstanceStep},
		{Name: "IsolateDBInstances", Tool: "IsolateDBInstances", Prepare: prepareIsolateStep},
		{Name: "DisIsolateDBInstances", Tool: "DisIsolateDBInstances", Prepare: prepareDisIsolateStep},
		{Name: "ModifyDBInstanceParameters", Tool: "ModifyDBInstanceParameters", Prepare: prepareInstanceStep, Verify: verifyParameterChange},
		{Name: "CreateInstances", Tool: "CreateInstances", Prepare: prepareCreateInstancesStep},
		{Name: "ModifyDBInstanceSpec", Tool: "ModifyDBInstanceSpec", Prepare: prepareInstanceStep},
		{Name: "CloneDBInstance", Tool: "CloneDBInstance", Prepare: prepareCloneStep},
		{Name: "CreateReadOnlyDBInstance", Tool: "CreateReadOnlyDBInstance", Prepare: prepareCreateReadOnlyStep},
	}
}

func (r *Runner) Run(steps []StepRunner, only string) error {
	if only != "" {
		for _, step := range steps {
			if step.Name == only {
				_, err := r.runStep(step)
				return err
			}
		}
		return fmt.Errorf("step %q not found", only)
	}

	for _, step := range steps {
		if stop, err := r.runStep(step); err != nil {
			return err
		} else if stop {
			break
		}
	}
	return nil
}

func (r *Runner) runStep(step StepRunner) (bool, error) {
	cfg := r.plan.Step(step.Name)
	fmt.Printf("--- %s ---\n", step.Name)

	if !cfg.Enabled {
		r.skip(step.Name, "配置中 enabled=false")
		fmt.Println()
		return false, nil
	}
	if !cfg.Approved {
		r.skip(step.Name, cfg.approvalHint())
		fmt.Println()
		return false, nil
	}

	prepare := step.Prepare
	if prepare == nil {
		prepare = prepareInstanceStep
	}
	args, skipReason, err := prepare(r, cfg)
	if err != nil {
		r.fail(step.Name, fmt.Sprintf("prepare failed: %v", err))
		fmt.Println()
		return r.plan.Options.StopOnFailure, nil
	}
	if skipReason != "" {
		r.skip(step.Name, skipReason)
		fmt.Println()
		return false, nil
	}

	if cfg.Approved {
		args["confirm"] = true
	}

	out, err := r.call(step.Tool, args)
	ok := false
	if err == nil && isIdempotentConflict(step.Name, out) {
		show := out
		if len(show) > 320 {
			show = show[:320] + "...(truncated)"
		}
		fmt.Printf("~ %s: %s\n", step.Name, show)
		r.passed++
		ok = true
	} else {
		ok = r.check(step.Name, out, err)
	}
	if !ok {
		fmt.Println()
		return r.plan.Options.StopOnFailure, nil
	}

	if shouldWaitForRunning(step.Name) && r.targetsPrimaryInstance(args) {
		if err := r.waitForTargetRunning(step.Name); err != nil {
			r.fail(step.Name+" wait", err.Error())
			fmt.Println()
			return r.plan.Options.StopOnFailure, nil
		}
	}

	if step.Verify != nil {
		if err := step.Verify(r, cfg, args, out); err != nil {
			r.fail(step.Name+" verify", err.Error())
			fmt.Println()
			return r.plan.Options.StopOnFailure, nil
		}
	}

	pause := cfg.pauseSeconds(r.plan.Options.DefaultPauseSeconds)
	if pause > 0 {
		time.Sleep(time.Duration(pause) * time.Second)
	}

	fmt.Println()
	return false, nil
}

func (r *Runner) call(name string, args map[string]interface{}) (string, error) {
	req := mcp.CallToolRequest{}
	req.Params.Name = "postgres-" + name
	req.Params.Arguments = args

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	res, err := r.client.CallTool(ctx, req)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, content := range res.Content {
		if tc, ok := content.(mcp.TextContent); ok {
			sb.WriteString(tc.Text)
		}
	}
	return sb.String(), nil
}

func (r *Runner) check(label, out string, err error) bool {
	if err != nil {
		r.fail(label, err.Error())
		return false
	}
	var raw map[string]interface{}
	if json.Unmarshal([]byte(out), &raw) == nil {
		if code, ok := raw["code"].(float64); ok && code != 0 && code != 200 {
			r.fail(label, out)
			return false
		}
	}
	show := out
	if len(show) > 320 {
		show = show[:320] + "...(truncated)"
	}
	fmt.Printf("✓ %s: %s\n", label, show)
	r.passed++
	return true
}

func (r *Runner) fail(label, reason string) {
	fmt.Printf("✗ %s: %s\n", label, reason)
	r.failed++
}

func (r *Runner) skip(label, reason string) {
	fmt.Printf("- %s [SKIPPED]: %s\n", label, reason)
	r.skipped++
}

func prepareInstanceStep(r *Runner, cfg StepConfig) (map[string]interface{}, string, error) {
	args := cloneMap(cfg.Args)
	args["region"] = r.plan.Target.Region
	if _, ok := args["DBInstanceId"]; !ok {
		args["DBInstanceId"] = r.plan.Target.DBInstanceID
	}
	return args, "", nil
}

func prepareRegionOnlyStep(r *Runner, cfg StepConfig) (map[string]interface{}, string, error) {
	args := cloneMap(cfg.Args)
	args["region"] = r.plan.Target.Region
	return args, "", nil
}

func prepareCreateInstancesStep(r *Runner, cfg StepConfig) (map[string]interface{}, string, error) {
	args, _, err := prepareRegionOnlyStep(r, cfg)
	if err != nil {
		return nil, "", err
	}
	if err := r.applyCreateDefaults(args, true); err != nil {
		return nil, "", err
	}
	return args, "", nil
}

func prepareIsolateStep(r *Runner, cfg StepConfig) (map[string]interface{}, string, error) {
	args := cloneMap(cfg.Args)
	args["region"] = r.plan.Target.Region
	if _, ok := args["DBInstanceId"]; !ok {
		if _, setOK := args["DBInstanceIdSet"]; !setOK {
			args["DBInstanceId"] = r.plan.Target.DBInstanceID
		}
	}
	return args, "", nil
}

func prepareDisIsolateStep(r *Runner, cfg StepConfig) (map[string]interface{}, string, error) {
	return prepareIsolateStep(r, cfg)
}

func prepareCreateReadOnlyStep(r *Runner, cfg StepConfig) (map[string]interface{}, string, error) {
	args := cloneMap(cfg.Args)
	args["region"] = r.plan.Target.Region
	if _, ok := args["MasterDBInstanceId"]; !ok {
		if _, aliasOK := args["DBInstanceId"]; !aliasOK {
			args["MasterDBInstanceId"] = r.plan.Target.DBInstanceID
		}
	}
	if err := r.applyCreateDefaults(args, true); err != nil {
		return nil, "", err
	}
	return args, "", nil
}

func prepareCloneStep(r *Runner, cfg StepConfig) (map[string]interface{}, string, error) {
	args := cloneMap(cfg.Args)
	args["region"] = r.plan.Target.Region
	if _, ok := args["DBInstanceId"]; !ok {
		args["DBInstanceId"] = r.plan.Target.DBInstanceID
	}
	if err := r.applyCreateDefaults(args, true); err != nil {
		return nil, "", err
	}
	if !hasDBNodeSet(args["DBNodeSet"]) {
		meta, err := r.targetInstanceMeta(false)
		if err != nil {
			return nil, "", err
		}
		args["DBNodeSet"] = defaultDBNodeSet(meta.Zone)
	}
	if _, ok := args["AutoRenewFlag"]; !ok {
		args["AutoRenewFlag"] = 0
	}
	if stringValue(args["BackupSetId"]) == "" && stringValue(args["RecoveryTargetTime"]) == "" {
		backupID, err := r.findFirstBaseBackupID()
		if err != nil {
			return nil, "", err
		}
		if backupID == "" {
			return nil, "无可用基础备份，无法自动补全 BackupSetId", nil
		}
		args["BackupSetId"] = backupID
	}
	return args, "", nil
}

func prepareBackupDownloadURL(r *Runner, cfg StepConfig) (map[string]interface{}, string, error) {
	args := cloneMap(cfg.Args)
	args["region"] = r.plan.Target.Region
	if _, ok := args["DBInstanceId"]; !ok {
		args["DBInstanceId"] = r.plan.Target.DBInstanceID
	}
	if stringValue(args["BackupType"]) == "" {
		args["BackupType"] = "BaseBackup"
	}
	if stringValue(args["BackupId"]) == "" {
		backupID, err := r.findFirstBaseBackupID()
		if err != nil {
			return nil, "", err
		}
		if backupID == "" {
			return nil, "无可用基础备份，无法获取下载链接", nil
		}
		args["BackupId"] = backupID
	}
	return args, "", nil
}

func prepareModifySecurityGroups(r *Runner, cfg StepConfig) (map[string]interface{}, string, error) {
	args := cloneMap(cfg.Args)
	args["region"] = r.plan.Target.Region
	if _, ok := args["DBInstanceId"]; !ok {
		args["DBInstanceId"] = r.plan.Target.DBInstanceID
	}
	reuseCurrent := boolValue(args["ReuseCurrentSecurityGroupSet"])
	delete(args, "ReuseCurrentSecurityGroupSet")

	current, err := r.describeCurrentSecurityGroups()
	if err != nil {
		return nil, "", err
	}
	if len(current) == 0 {
		return nil, "当前实例未查到任何安全组，无法做全量替换测试", nil
	}
	if len(stringSliceValue(args["SecurityGroupIdSet"])) == 0 && len(stringSliceValue(args["SecurityGroupIds"])) == 0 {
		reuseCurrent = true
	}
	if reuseCurrent {
		args["SecurityGroupIdSet"] = current
	}
	return args, "", nil
}

func verifyInstanceName(r *Runner, _ StepConfig, args map[string]interface{}, _ string) error {
	want := stringValue(args["InstanceName"])
	if want == "" {
		return nil
	}
	deadline := time.Now().Add(30 * time.Second)
	lastGot := ""
	for {
		meta, err := r.targetInstanceMeta(true)
		if err == nil && meta != nil {
			lastGot = meta.Name
			if meta.Name == want {
				fmt.Printf("  verify instance name: %q\n", meta.Name)
				return nil
			}
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("instance name verify failed: want=%q got=%q", want, lastGot)
		}
		time.Sleep(3 * time.Second)
	}
}

func verifyAccountCreated(r *Runner, _ StepConfig, args map[string]interface{}, _ string) error {
	user := stringValue(args["UserName"])
	if user == "" {
		return nil
	}
	accounts, err := r.describeAccounts()
	if err != nil {
		return err
	}
	if !contains(accounts, user) {
		return fmt.Errorf("account %q not found after CreateAccount", user)
	}
	fmt.Printf("  verify account exists: %q\n", user)
	return nil
}

func verifyAccountDeleted(r *Runner, _ StepConfig, args map[string]interface{}, _ string) error {
	user := stringValue(args["UserName"])
	if user == "" {
		return nil
	}
	accounts, err := r.describeAccounts()
	if err != nil {
		return err
	}
	if contains(accounts, user) {
		return fmt.Errorf("account %q still exists after DeleteAccount", user)
	}
	fmt.Printf("  verify account deleted: %q\n", user)
	return nil
}

func verifyPrivilegeGrant(r *Runner, cfg StepConfig, args map[string]interface{}, _ string) error {
	privs, err := r.queryPrivileges(cfg, args)
	if err != nil {
		return err
	}
	expected := extractExpectedPrivileges(args)
	missing := diffStrings(expected, privs)
	if len(missing) > 0 {
		return fmt.Errorf("grant verify failed, missing privileges: %v (current=%v)", missing, privs)
	}
	fmt.Printf("  verify granted privileges: %v\n", privs)
	return nil
}

func verifyPrivilegeRevoke(r *Runner, cfg StepConfig, args map[string]interface{}, _ string) error {
	privs, err := r.queryPrivileges(cfg, args)
	if err != nil {
		return err
	}
	expected := extractExpectedPrivileges(args)
	stillPresent := intersectStrings(expected, privs)
	if len(stillPresent) > 0 {
		return fmt.Errorf("revoke verify failed, privileges still present: %v (current=%v)", stillPresent, privs)
	}
	fmt.Printf("  verify revoked privileges, current: %v\n", privs)
	return nil
}

func verifyDatabaseCreated(r *Runner, _ StepConfig, args map[string]interface{}, _ string) error {
	dbName := stringValue(args["DatabaseName"])
	if dbName == "" {
		return nil
	}
	owner, found, err := r.findDatabaseOwner(dbName)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("database %q not found after CreateDatabase", dbName)
	}
	fmt.Printf("  verify database exists: %q owner=%q\n", dbName, owner)
	return nil
}

func verifyDatabaseOwner(r *Runner, _ StepConfig, args map[string]interface{}, _ string) error {
	dbName := stringValue(args["DatabaseName"])
	ownerWant := stringValue(args["DatabaseOwner"])
	if dbName == "" || ownerWant == "" {
		return nil
	}
	owner, found, err := r.findDatabaseOwner(dbName)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("database %q not found when verifying owner", dbName)
	}
	if owner != ownerWant {
		return fmt.Errorf("database owner verify failed: want=%q got=%q", ownerWant, owner)
	}
	fmt.Printf("  verify database owner: %q -> %q\n", dbName, owner)
	return nil
}

func verifySecurityGroups(r *Runner, _ StepConfig, args map[string]interface{}, _ string) error {
	want := stringSliceValue(args["SecurityGroupIdSet"])
	if len(want) == 0 {
		want = stringSliceValue(args["SecurityGroupIds"])
	}
	if len(want) == 0 {
		return nil
	}
	got, err := r.describeCurrentSecurityGroups()
	if err != nil {
		return err
	}
	sort.Strings(want)
	sort.Strings(got)
	if strings.Join(want, ",") != strings.Join(got, ",") {
		return fmt.Errorf("security groups verify failed: want=%v got=%v", want, got)
	}
	fmt.Printf("  verify security groups: %v\n", got)
	return nil
}

func verifyParameterChange(r *Runner, _ StepConfig, args map[string]interface{}, _ string) error {
	rawList, ok := args["ParamList"].([]interface{})
	if !ok || len(rawList) == 0 {
		return nil
	}
	first, ok := rawList[0].(map[string]interface{})
	if !ok {
		return nil
	}
	name := stringValue(first["Name"])
	want := stringValue(first["ExpectedValue"])
	if want == "" {
		want = stringValue(first["Value"])
	}
	if name == "" || want == "" {
		return nil
	}
	out, err := r.call("DescribeDBInstanceParameters", map[string]interface{}{
		"region":       r.plan.Target.Region,
		"DBInstanceId": r.plan.Target.DBInstanceID,
		"ParamName":    name,
	})
	if err != nil {
		return err
	}
	got := extractParamCurrentValue(out, name)
	if got != want {
		return fmt.Errorf("parameter verify failed: %s want=%q got=%q", name, want, got)
	}
	fmt.Printf("  verify parameter %s=%s\n", name, got)
	return nil
}

func (r *Runner) describeAccounts() ([]string, error) {
	out, err := r.call("DescribeAccounts", map[string]interface{}{
		"region":       r.plan.Target.Region,
		"DBInstanceId": r.plan.Target.DBInstanceID,
	})
	if err != nil {
		return nil, err
	}
	return extractAccounts(out), nil
}

func (r *Runner) describeCurrentSecurityGroups() ([]string, error) {
	out, err := r.call("DescribeDBInstanceSecurityGroups", map[string]interface{}{
		"region":       r.plan.Target.Region,
		"DBInstanceId": r.plan.Target.DBInstanceID,
	})
	if err != nil {
		return nil, err
	}
	return extractSecurityGroupIDs(out), nil
}

func (r *Runner) findFirstBaseBackupID() (string, error) {
	out, err := r.call("DescribeBaseBackups", map[string]interface{}{
		"region":       r.plan.Target.Region,
		"DBInstanceId": r.plan.Target.DBInstanceID,
	})
	if err != nil {
		return "", err
	}
	return extractFirstBackupID(out), nil
}

func (r *Runner) findDatabaseOwner(dbName string) (string, bool, error) {
	out, err := r.call("DescribeDatabases", map[string]interface{}{
		"region":       r.plan.Target.Region,
		"DBInstanceId": r.plan.Target.DBInstanceID,
	})
	if err != nil {
		return "", false, err
	}
	owner, found := extractDatabaseOwner(out, dbName)
	return owner, found, nil
}

func (r *Runner) targetInstanceMeta(refresh bool) (*InstanceMeta, error) {
	if !refresh && r.targetMeta != nil {
		meta := *r.targetMeta
		return &meta, nil
	}
	out, err := r.call("DescribeDBInstances", map[string]interface{}{
		"region": r.plan.Target.Region,
		"Filters": []map[string]interface{}{
			{
				"Name":   "db-instance-id",
				"Values": []string{r.plan.Target.DBInstanceID},
			},
		},
		"Limit": 20,
	})
	if err != nil {
		return nil, err
	}
	meta, found := extractInstanceMeta(out, r.plan.Target.DBInstanceID)
	if !found {
		return nil, fmt.Errorf("target instance %q not found in DescribeDBInstances", r.plan.Target.DBInstanceID)
	}
	r.targetMeta = &meta
	metaCopy := meta
	return &metaCopy, nil
}

func (r *Runner) applyCreateDefaults(args map[string]interface{}, includeAutoRenew bool) error {
	meta, err := r.targetInstanceMeta(false)
	if err != nil {
		return err
	}
	if stringValue(args["SpecCode"]) == "" {
		args["SpecCode"] = meta.SpecCode
	}
	if stringValue(args["VpcId"]) == "" {
		args["VpcId"] = meta.VpcId
	}
	if stringValue(args["SubnetId"]) == "" {
		args["SubnetId"] = meta.SubnetId
	}
	if includeAutoRenew {
		if _, ok := args["AutoRenewFlag"]; !ok {
			args["AutoRenewFlag"] = 0
		}
	}
	return nil
}

func hasDBNodeSet(raw interface{}) bool {
	switch items := raw.(type) {
	case []interface{}:
		return len(items) > 0
	case []map[string]interface{}:
		return len(items) > 0
	case []map[string]string:
		return len(items) > 0
	default:
		return false
	}
}

func defaultDBNodeSet(zone string) []map[string]interface{} {
	if zone == "" {
		return nil
	}
	return []map[string]interface{}{
		{"Role": "Primary", "Zone": zone},
		{"Role": "Standby", "Zone": zone},
	}
}

func (r *Runner) targetsPrimaryInstance(args map[string]interface{}) bool {
	if stringValue(args["DBInstanceId"]) == r.plan.Target.DBInstanceID {
		return true
	}
	if stringValue(args["MasterDBInstanceId"]) == r.plan.Target.DBInstanceID {
		return true
	}
	for _, id := range stringSliceValue(args["DBInstanceIdSet"]) {
		if id == r.plan.Target.DBInstanceID {
			return true
		}
	}
	return false
}

func (r *Runner) waitForTargetRunning(stepName string) error {
	deadline := time.Now().Add(waitTimeout(r.timeout))
	lastStatus := ""
	for {
		meta, err := r.targetInstanceMeta(true)
		if err == nil && meta != nil {
			lastStatus = meta.Status
			if strings.EqualFold(meta.Status, "running") {
				fmt.Printf("  wait %s: instance status=%s\n", stepName, meta.Status)
				return nil
			}
		}
		if time.Now().After(deadline) {
			if lastStatus == "" {
				lastStatus = "unknown"
			}
			return fmt.Errorf("wait target instance running timeout after %s, last status=%q", waitTimeout(r.timeout), lastStatus)
		}
		time.Sleep(5 * time.Second)
	}
}

func waitTimeout(base time.Duration) time.Duration {
	if base <= 0 {
		return 5 * time.Minute
	}
	wait := base * 5
	if wait < 5*time.Minute {
		return 5 * time.Minute
	}
	if wait > 15*time.Minute {
		return 15 * time.Minute
	}
	return wait
}

func isIdempotentConflict(stepName, out string) bool {
	lower := strings.ToLower(out)
	switch stepName {
	case "CreateAccount":
		return strings.Contains(lower, "account already exist")
	case "CreateDatabase":
		return strings.Contains(lower, "database already exists")
	default:
		return false
	}
}

func shouldWaitForRunning(stepName string) bool {
	switch stepName {
	case "IsolateDBInstances":
		return false
	default:
		return true
	}
}

func (r *Runner) queryPrivileges(cfg StepConfig, args map[string]interface{}) ([]string, error) {
	verifyArgs := cloneMap(cfg.VerifyArgs)
	if len(verifyArgs) == 0 {
		verifyArgs = derivePrivilegeVerifyArgs(args)
	}
	verifyArgs["region"] = r.plan.Target.Region
	if _, ok := verifyArgs["DBInstanceId"]; !ok {
		verifyArgs["DBInstanceId"] = r.plan.Target.DBInstanceID
	}
	if _, ok := verifyArgs["UserName"]; !ok {
		verifyArgs["UserName"] = args["UserName"]
	}
	if len(stringSliceValue(verifyArgs["DatabaseObjectSet"])) == 0 {
		if raw, ok := verifyArgs["DatabaseObjectSet"].([]interface{}); !ok || len(raw) == 0 {
			return nil, fmt.Errorf("verify_args.DatabaseObjectSet is required for privilege verification")
		}
	}
	out, err := r.call("DescribeAccountPrivileges", verifyArgs)
	if err != nil {
		return nil, err
	}
	return extractPrivilegeSet(out), nil
}

func derivePrivilegeVerifyArgs(args map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	result["UserName"] = args["UserName"]
	rawSet, ok := args["ModifyPrivilegeSet"].([]interface{})
	if !ok || len(rawSet) == 0 {
		return result
	}
	entry, ok := rawSet[0].(map[string]interface{})
	if !ok {
		return result
	}
	dbPriv, ok := entry["DatabasePrivilege"].(map[string]interface{})
	if !ok {
		return result
	}
	obj, ok := dbPriv["Object"].(map[string]interface{})
	if !ok {
		return result
	}
	result["DatabaseObjectSet"] = []map[string]interface{}{obj}
	return result
}

func extractExpectedPrivileges(args map[string]interface{}) []string {
	rawSet, ok := args["ModifyPrivilegeSet"].([]interface{})
	if !ok || len(rawSet) == 0 {
		return nil
	}
	entry, ok := rawSet[0].(map[string]interface{})
	if !ok {
		return nil
	}
	dbPriv, ok := entry["DatabasePrivilege"].(map[string]interface{})
	if !ok {
		return nil
	}
	return stringSliceValue(dbPriv["PrivilegeSet"])
}

func extractAccounts(out string) []string {
	resp := responseObject(out)
	for _, v := range resp {
		arr, ok := v.([]interface{})
		if !ok {
			continue
		}
		var result []string
		for _, item := range arr {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if name := stringValue(m["UserName"]); name != "" {
				result = append(result, name)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return nil
}

func extractSecurityGroupIDs(out string) []string {
	resp := responseObject(out)
	rawSet, ok := resp["SecurityGroupSet"].([]interface{})
	if !ok {
		return nil
	}
	var ids []string
	for _, item := range rawSet {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if id := stringValue(m["SecurityGroupId"]); id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

func extractFirstBackupID(out string) string {
	resp := responseObject(out)
	for _, v := range resp {
		arr, ok := v.([]interface{})
		if !ok || len(arr) == 0 {
			continue
		}
		first, ok := arr[0].(map[string]interface{})
		if !ok {
			continue
		}
		if id := stringValue(first["Id"]); id != "" {
			return id
		}
	}
	return ""
}

func extractDatabaseOwner(out, dbName string) (string, bool) {
	resp := responseObject(out)
	for _, v := range resp {
		arr, ok := v.([]interface{})
		if !ok {
			continue
		}
		for _, item := range arr {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if stringValue(m["DatabaseName"]) == dbName {
				return stringValue(m["DatabaseOwner"]), true
			}
		}
	}
	return "", false
}

func extractInstanceMeta(out, targetID string) (InstanceMeta, bool) {
	resp := responseObject(out)
	rawSet, ok := resp["DBInstanceSet"].([]interface{})
	if !ok {
		return InstanceMeta{}, false
	}
	for _, item := range rawSet {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		id := stringValue(m["DBInstanceId"])
		if targetID != "" && id != targetID {
			continue
		}
		return InstanceMeta{
			ID:             id,
			Name:           stringValue(m["DBInstanceName"]),
			Status:         stringValue(m["DBInstanceStatus"]),
			Zone:           stringValue(m["Zone"]),
			VpcId:          stringValue(m["VpcId"]),
			SubnetId:       stringValue(m["SubnetId"]),
			SpecCode:       stringValue(m["DBInstanceClass"]),
			DBMajorVersion: stringValue(m["DBMajorVersion"]),
		}, true
	}
	return InstanceMeta{}, false
}

func extractPrivilegeSet(out string) []string {
	resp := responseObject(out)
	rawSet, ok := resp["PrivilegeSet"].([]interface{})
	if !ok || len(rawSet) == 0 {
		return nil
	}
	first, ok := rawSet[0].(map[string]interface{})
	if !ok {
		return nil
	}
	return stringSliceValue(first["PrivilegeSet"])
}

func extractParamCurrentValue(out, name string) string {
	resp := responseObject(out)
	rawList, ok := resp["Detail"].([]interface{})
	if !ok {
		return ""
	}
	for _, item := range rawList {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if stringValue(m["Name"]) == name {
			return stringValue(m["CurrentValue"])
		}
	}
	return ""
}

func extractTopStringField(out, field string) string {
	resp := responseObject(out)
	return stringValue(resp[field])
}

func responseObject(out string) map[string]interface{} {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return map[string]interface{}{}
	}
	resp, ok := raw["Response"].(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
	}
	return resp
}

func cloneMap(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return map[string]interface{}{}
	}
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func stringValue(v interface{}) string {
	switch x := v.(type) {
	case string:
		return x
	case fmt.Stringer:
		return x.String()
	case int:
		return fmt.Sprintf("%d", x)
	case int64:
		return fmt.Sprintf("%d", x)
	case float64:
		if x == float64(int64(x)) {
			return fmt.Sprintf("%d", int64(x))
		}
		return fmt.Sprintf("%v", x)
	default:
		return ""
	}
}

func boolValue(v interface{}) bool {
	switch x := v.(type) {
	case bool:
		return x
	case string:
		return strings.EqualFold(x, "true") || x == "1"
	case int:
		return x != 0
	case int64:
		return x != 0
	case float64:
		return x != 0
	default:
		return false
	}
}

func stringSliceValue(v interface{}) []string {
	switch x := v.(type) {
	case []string:
		return append([]string(nil), x...)
	case []interface{}:
		result := make([]string, 0, len(x))
		for _, item := range x {
			if s := stringValue(item); s != "" {
				result = append(result, s)
			}
		}
		return result
	default:
		return nil
	}
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func diffStrings(expected, actual []string) []string {
	var missing []string
	for _, item := range expected {
		if !contains(actual, item) {
			missing = append(missing, item)
		}
	}
	return missing
}

func intersectStrings(expected, actual []string) []string {
	var shared []string
	for _, item := range expected {
		if contains(actual, item) {
			shared = append(shared, item)
		}
	}
	return shared
}
