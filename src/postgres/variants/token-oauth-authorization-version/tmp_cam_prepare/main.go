package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

func parseEnv(path string) map[string]string {
	out := map[string]string{}
	data, err := os.ReadFile(path)
	if err != nil {
		return out
	}
	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		if len(v) >= 2 && ((v[0] == '\'' && v[len(v)-1] == '\'') || (v[0] == '"' && v[len(v)-1] == '"')) {
			v = v[1 : len(v)-1]
		}
		out[k] = v
	}
	return out
}

func str(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func mustCred() *common.Credential {
	env := parseEnv("/Users/lumenyang/workspace/tencentcloud-mcp-server/src/postgres/.env")
	sid := env["MCP_SECRET_ID"]
	if sid == "" {
		sid = env["TENCENTCLOUD_SECRET_ID"]
	}
	sk := env["MCP_SECRET_KEY"]
	if sk == "" {
		sk = env["TENCENTCLOUD_SECRET_KEY"]
	}
	tok := env["MCP_SESSION_TOKEN"]
	if tok == "" {
		tok = env["TENCENTCLOUD_SESSION_TOKEN"]
	}
	if sid == "" || sk == "" {
		fmt.Println("CREDENTIAL_MISSING")
		os.Exit(1)
	}
	if tok != "" {
		return common.NewTokenCredential(sid, sk, tok)
	}
	return common.NewCredential(sid, sk)
}

func main() {
	cred := mustCred()
	cpf := profile.NewClientProfile()

	accountID := "100023351712"
	userID := "100050127078"
	callerArn := fmt.Sprintf("qcs::cam::uin/%s:uin/%s", accountID, userID)
	fmt.Println("CALLER_ARN=" + callerArn)
	fmt.Println("ACCOUNT_ID=" + accountID)
	fmt.Println("USER_ID=" + userID)

	roleName := fmt.Sprintf("pg-mcp-verify-%d", time.Now().Unix())
	trustPolicy := map[string]any{
		"version": "2.0",
		"statement": []map[string]any{{
			"effect": "allow",
			"action": "name/sts:AssumeRole",
			"principal": map[string]any{
				"qcs": []string{fmt.Sprintf("qcs::cam::uin/%s:root", accountID)},
			},
			"condition": map[string]any{
				"string_equal": map[string]any{
					"qcs:assume_principal_arn": []string{callerArn},
				},
			},
		}},
	}
	trustJSON, _ := json.Marshal(trustPolicy)

	camClient, err := cam.NewClient(cred, "ap-guangzhou", cpf)
	if err != nil {
		fmt.Println("CAM_CLIENT_ERR=" + err.Error())
		os.Exit(1)
	}
	createReq := cam.NewCreateRoleRequest()
	createReq.RoleName = common.StringPtr(roleName)
	createReq.PolicyDocument = common.StringPtr(string(trustJSON))
	desc := "Temporary role for pg-mcp assume-role validation"
	createReq.Description = common.StringPtr(desc)
	consoleLogin := uint64(0)
	createReq.ConsoleLogin = &consoleLogin
	sessionDuration := uint64(1800)
	createReq.SessionDuration = &sessionDuration
	_, err = camClient.CreateRole(createReq)
	if err != nil {
		fmt.Println("CREATE_ROLE_ERR=" + err.Error())
		os.Exit(1)
	}

	attachReq := cam.NewAttachRolePolicyRequest()
	attachReq.AttachRoleName = common.StringPtr(roleName)
	policyName := "QcloudPostgreSQLReadOnlyAccess"
	attachReq.PolicyName = common.StringPtr(policyName)
	_, err = camClient.AttachRolePolicy(attachReq)
	if err != nil {
		fmt.Println("ATTACH_POLICY_ERR=" + err.Error())
		fmt.Printf("ROLE_NAME=%s\n", roleName)
		os.Exit(1)
	}

	fmt.Printf("ROLE_NAME=%s\n", roleName)
	fmt.Printf("ROLE_ARN=qcs::cam::uin/%s:roleName/%s\n", accountID, roleName)
	fmt.Println("POLICY_NAME=QcloudPostgreSQLReadOnlyAccess")
}
