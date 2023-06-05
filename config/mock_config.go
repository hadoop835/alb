package config

import (
	"fmt"
	"time"
)

// a mock of config impl IConfig
type MockConfig struct {
	Name                  string
	NS                    string
	Domain                string
	Pod                   string
	LeaderLabel           string
	Http                  int
	Https                 int
	SSlCert               string
	SSLStrategy           string
	LeaderConfig          LeaderConfig
	GatewayCfg            GatewayCfg
	Mode                  Mode
	ControllerNetWorkMode ControllerNetWorkMode
	EnableAlb             bool
}

func DefaultMock() *MockConfig {
	return &MockConfig{
		Name:                  "alb-dev",
		NS:                    "cpaas-system",
		Domain:                "cpaas.io",
		Pod:                   "mock-pod",
		LeaderLabel:           "alb2.cpaas.ip/leader",
		Http:                  80,
		Https:                 443,
		Mode:                  Controller,
		ControllerNetWorkMode: Host,
		LeaderConfig: LeaderConfig{
			LeaseDuration: time.Second * 15,
			RenewDeadline: time.Second * 10,
			RetryPeriod:   time.Second * 2,
		},
	}
}

func (c *MockConfig) GetNs() string {
	return c.NS
}

func (c *MockConfig) GetAlbName() string {
	return c.Name
}

func (c *MockConfig) GetDomain() string {
	return c.Domain
}

func (c *MockConfig) GetPodName() string {
	return c.Pod
}

func (c *MockConfig) GetLabelLeader() string {
	return c.LeaderLabel
}

func (c *MockConfig) GetLabelAlbName() string {
	return "alb-name"
}
func (c *MockConfig) GetLabelSourceType() string {
	return "source"
}

func (c *MockConfig) GetIngressHttpPort() int {
	return c.Http
}

func (c *MockConfig) GetIngressHttpsPort() int {
	return c.Https
}

func (c *MockConfig) GetLabelSourceIngressVer() string {
	return "ingress-ver"
}

func (c *MockConfig) GetDefaultSSLSCert() string {
	return c.SSlCert
}

func (c *MockConfig) GetDefaultSSLStrategy() string {
	return c.SSLStrategy
}

func (c *MockConfig) GetLeaderConfig() LeaderConfig {
	return c.LeaderConfig
}

func (c *MockConfig) DebugRuleSync() bool {
	return true
}
func (c *MockConfig) GetLabelSourceIngressPathIndex() string {
	return "path-index"
}

func (c *MockConfig) GetLabelSourceIngressRuleIndex() string {
	return "rule-index"
}

func (c *MockConfig) GetGatewayCfg() GatewayCfg {
	return c.GatewayCfg
}

func (c *MockConfig) GetMode() Mode {
	return c.Mode
}

func (c *MockConfig) GetNetworkMode() ControllerNetWorkMode {
	return c.ControllerNetWorkMode
}

func (c *MockConfig) IsEnableAlb() bool {
	return c.EnableAlb
}

func (c *MockConfig) GetAnnotationIngressAddress() string {
	return fmt.Sprintf(INGRESS_ADDRESS_NAME, c.GetDomain(), c.GetAlbName())
}

func (c *MockConfig) IsEnableVIP() bool {
	return false
}
