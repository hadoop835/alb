package controller

import (
	"encoding/json"
	"fmt"

	"alauda.io/alb2/config"
	. "alauda.io/alb2/controller/types"
	"k8s.io/klog/v2"
)

func GetAlbIngressRewriteResponseAnnotation() string {
	return fmt.Sprintf("alb.ingress.%s/rewrite-response", config.Get("DOMAIN"))
}
func GetAlbRuleRewriteResponseAnnotation() string {
	return fmt.Sprintf("alb.rule.%s/rewrite-response", config.Get("DOMAIN"))
}

func rewriteResponseConfigFromJson(jsonStr string) (*RewriteResponseConfig, error) {
	cfg := RewriteResponseConfig{}
	err := json.Unmarshal([]byte(jsonStr), &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.IsEmpty() {
		return nil, fmt.Errorf("empty config")
	}
	return &cfg, err
}

func GenerateRuleAnnotationFromIngressAnnotation(ingressName string, annotation map[string]string) map[string]string {

	ruleAnnotation := make(map[string]string)

	if val, ok := annotation[GetAlbIngressRewriteResponseAnnotation()]; ok {
		_, err := rewriteResponseConfigFromJson(val)
		if err != nil {
			klog.Errorf("ext ingress rewrite_response: invalid annotation in ingress '%v' annotation is '%v' err %v", ingressName, val, err)
		} else {
			ruleAnnotation[GetAlbRuleRewriteResponseAnnotation()] = val
		}
	}
	return ruleAnnotation
}

func RuleConfigFromRuleAnnotation(ruleName string, annotation map[string]string) *RuleConfig {
	cfg := RuleConfig{}

	if val, ok := annotation[GetAlbRuleRewriteResponseAnnotation()]; ok {
		rewriteCfg, err := rewriteResponseConfigFromJson(val)
		if err != nil {
			klog.Errorf("ext rule rewrite_response: invalid annotation in rule '%v' annotation is '%v' err %v", ruleName, val, err)
		} else {
			cfg.RewriteResponse = rewriteCfg
		}
	}
	if cfg.IsEmpty() {
		return nil
	}
	return &cfg
}
