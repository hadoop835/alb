package controllers

import (
	albv1 "alauda.io/alb2/pkg/apis/alauda/v1"
	albv2 "alauda.io/alb2/pkg/apis/alauda/v2beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	gv1a2t "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gv1b1t "sigs.k8s.io/gateway-api/apis/v1beta1"
)

func InitScheme(scheme *runtime.Scheme) *runtime.Scheme {
	_ = albv2.AddToScheme(scheme)
	_ = appsv1.AddToScheme(scheme)
	_ = albv1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)
	_ = netv1.AddToScheme(scheme)
	_ = gv1b1t.AddToScheme(scheme)
	_ = gv1a2t.AddToScheme(scheme)
	return scheme
}
