package ingress

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/cache"
	"strings"

	"alauda.io/alb2/config"
	m "alauda.io/alb2/modules"
	"github.com/fatih/set"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/klog/v2"
)

func parseSSLAnnotation(sslAnno string) map[string]string {
	// alb.networking.{domain}/tls: qq.com=cpaas-system/dex.tls,qq1.com=cpaas-system/dex1.tls
	if sslAnno == "" {
		return nil
	}
	rv := make(map[string]string)
	parts := strings.Split(sslAnno, ",")
	for _, p := range parts {
		kv := strings.Split(strings.TrimSpace(p), "=")
		if len(kv) != 2 {
			klog.Warningf("invalid ssl annotation format, %s", p)
			return nil
		}
		k, v := kv[0], kv[1]
		if rv[k] != "" {
			klog.Warningf("invalid ssl annotation duplicate host, %s", p)
			return nil
		}
		// ["qq.com"]="cpaas-system/dex.tls"
		rv[k] = v
	}
	return rv
}

func isDefaultBackend(ing *networkingv1.Ingress) bool {
	return len(ing.Spec.Rules) == 0 &&
		ing.Spec.DefaultBackend != nil &&
		ing.Spec.DefaultBackend.Resource == nil &&
		ing.Spec.DefaultBackend.Service != nil
}

func getIngressFtTypes(ing *networkingv1.Ingress) set.Interface {
	defaultSSLStrategy := config.Get("DEFAULT-SSL-STRATEGY")

	ALBSSLStrategyAnnotation := fmt.Sprintf("alb.networking.%s/enable-https", config.Get("DOMAIN"))
	ALBSSLAnnotation := fmt.Sprintf("alb.networking.%s/tls", config.Get("DOMAIN"))

	ingSSLStrategy := ing.Annotations[ALBSSLStrategyAnnotation]
	sslMap := parseSSLAnnotation(ing.Annotations[ALBSSLAnnotation])
	certs := make(map[string]string)
	for host, cert := range sslMap {
		if certs[strings.ToLower(host)] == "" {
			certs[strings.ToLower(host)] = cert
		}
	}

	needFtTypes := set.New(set.NonThreadSafe)
	for _, r := range ing.Spec.Rules {
		foundTLS := false
		for _, tls := range ing.Spec.TLS {
			for _, host := range tls.Hosts {
				if strings.EqualFold(r.Host, host) {
					needFtTypes.Add(m.ProtoHTTPS)
					foundTLS = true
				}
			}
		}
		if certs[strings.ToLower(r.Host)] != "" {
			needFtTypes.Add(m.ProtoHTTPS)
			foundTLS = true
		}
		if foundTLS == false {
			if defaultSSLStrategy == BothSSLStrategy && ingSSLStrategy != "false" {
				needFtTypes.Add(m.ProtoHTTPS)
				needFtTypes.Add(m.ProtoHTTP)
			} else if (defaultSSLStrategy == AlwaysSSLStrategy && ingSSLStrategy != "false") ||
				(defaultSSLStrategy == RequestSSLStrategy && ingSSLStrategy == "true") {
				needFtTypes.Add(m.ProtoHTTPS)
			} else {
				needFtTypes.Add(m.ProtoHTTP)
			}
		}
	}
	return needFtTypes
}

func ToInStr(backendPort networkingv1.ServiceBackendPort) intstr.IntOrString {
	intStrType := intstr.Int
	if backendPort.Number == 0 {
		intStrType = intstr.String
	}
	return intstr.IntOrString{Type: intStrType, IntVal: backendPort.Number, StrVal: backendPort.Name}
}

type NotExistsError string

// Error implements the error interface.
func (e NotExistsError) Error() string {
	return fmt.Sprintf("no object matching key %q in local store", string(e))
}

type ingressClassLister struct {
	cache.Store
}

// ByKey returns the Ingress matching key in the local Ingress Store.
func (il ingressClassLister) ByKey(key string) (*networkingv1.IngressClass, error) {
	i, exists, err := il.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, NotExistsError(key)
	}
	return i.(*networkingv1.IngressClass), nil
}

func CheckIngressClass(ingclass *networkingv1.IngressClass, icConfig *config.IngressClassConfiguration) bool {
	foundClassByName := false
	if icConfig.IngressClassByName && ingclass.Name == icConfig.AnnotationValue {
		foundClassByName = true
	}
	if !foundClassByName && ingclass.Spec.Controller != icConfig.Controller {
		klog.InfoS("ignoring ingressClass as the spec.controller is not the same with alb2", "ingressClass", klog.KObj(ingclass))
		return false
	}
	return true
}

func (c *Controller) GetIngressClass(ing *networkingv1.Ingress, icConfig *config.IngressClassConfiguration) (string, error) {
	// First we try ingressClassName
	if !icConfig.IgnoreIngressClass && ing.Spec.IngressClassName != nil {
		iclass, err := c.ingressClassLister.ByKey(*ing.Spec.IngressClassName)
		if err != nil {
			return "", err
		}
		return iclass.Name, nil
	}

	// Then we try annotation
	if ingressClass, ok := ing.GetAnnotations()[config.IngressKey]; ok {
		if ingressClass != "" && ingressClass != config.Get("NAME") {
			return "", fmt.Errorf("invalid ingress class annotation: %s", ingressClass)
		}
		return ingressClass, nil
	}

	// Then we accept if the WithoutClass is enabled
	if icConfig.WatchWithoutClass {
		// Reserving "_" as a "wildcard" name
		return "_", nil
	}
	return "", fmt.Errorf("ingress does not contain a valid IngressClass")
}
