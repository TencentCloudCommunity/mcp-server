package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type TestPlan struct {
	Server  ServerConfig         `yaml:"server"`
	Target  TargetConfig         `yaml:"target"`
	Options OptionsConfig        `yaml:"options"`
	Steps   map[string]StepConfig `yaml:"steps"`
}

type ServerConfig struct {
	URL string `yaml:"url"`
}

type TargetConfig struct {
	Region       string `yaml:"region"`
	DBInstanceID string `yaml:"db_instance_id"`
}

type OptionsConfig struct {
	TimeoutSeconds     int  `yaml:"timeout_seconds"`
	DefaultPauseSeconds int `yaml:"default_pause_seconds"`
	StopOnFailure      bool `yaml:"stop_on_failure"`
}

type StepConfig struct {
	Enabled      bool                   `yaml:"enabled"`
	Approved     bool                   `yaml:"approved"`
	GateNote     string                 `yaml:"gate_note"`
	SkipReason   string                 `yaml:"skip_reason"`
	PauseSeconds int                    `yaml:"pause_seconds"`
	Args         map[string]interface{} `yaml:"args"`
	VerifyArgs   map[string]interface{} `yaml:"verify_args"`
}

func LoadPlan(path string) (*TestPlan, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	var plan TestPlan
	if err := yaml.Unmarshal(data, &plan); err != nil {
		return nil, fmt.Errorf("parse yaml failed: %w", err)
	}

	plan.normalize()
	if err := plan.validate(); err != nil {
		return nil, err
	}
	return &plan, nil
}

func (p *TestPlan) normalize() {
	if p.Server.URL == "" {
		p.Server.URL = "http://127.0.0.1:9000/sse"
	}
	if p.Options.TimeoutSeconds <= 0 {
		p.Options.TimeoutSeconds = 60
	}
	if p.Options.DefaultPauseSeconds < 0 {
		p.Options.DefaultPauseSeconds = 0
	}
	if p.Steps == nil {
		p.Steps = map[string]StepConfig{}
	}
}

func (p *TestPlan) validate() error {
	if p.Target.Region == "" {
		return fmt.Errorf("target.region is required")
	}
	if p.Target.DBInstanceID == "" {
		return fmt.Errorf("target.db_instance_id is required")
	}
	return nil
}

func (p *TestPlan) Step(name string) StepConfig {
	if p == nil || p.Steps == nil {
		return StepConfig{}
	}
	return p.Steps[name]
}

func (s StepConfig) pauseSeconds(defaultPause int) int {
	if s.PauseSeconds > 0 {
		return s.PauseSeconds
	}
	if s.PauseSeconds < 0 {
		return 0
	}
	return defaultPause
}

func (s StepConfig) approvalHint() string {
	if s.SkipReason != "" {
		return s.SkipReason
	}
	if s.GateNote != "" {
		return s.GateNote
	}
	return "未在配置文件中批准"
}
