package kong

import (
	"github.com/ciroque/k8s-kong-federated-ingress/internal/k8s"
	v1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

func compareStringArrays(l []string, r []string) bool {
	if len(l) != len(r) {
		return false
	}
	for i, e := range l {
		if r[i] != e {
			return false
		}
	}
	return true
}

func servicesMatch(l k8s.Service, r k8s.Service) bool {
	return l.Name == r.Name &&
		compareStringArrays(l.Paths, r.Paths) &&
		l.Port == r.Port &&
		compareStringArrays(l.Addresses, r.Addresses)
}

func TestTranslation_IngressToService(t *testing.T) {
	translation := new(Translation)
	testHost := "test-host"
	testNamespace := "testing-namespace"
	testServiceName := "test-service"
	addresses := []string{
		"10.200.30.400:80",
		"10.200.30.401:80",
	}
	expectedService := k8s.Service{
		Addresses: addresses,
		Name:      translation.FormatServiceName(testNamespace, testServiceName),
		Paths:     []string{"/apple", "/banana"},
		Port:      80,
	}

	ingress := networking.Ingress{
		Status: networking.IngressStatus{
			LoadBalancer: v1.LoadBalancerStatus{
				Ingress: []v1.LoadBalancerIngress{
					{
						IP: "10.200.30.400",
					},
					{
						IP: "10.200.30.401",
					},
				},
			},
		},
		Spec: networking.IngressSpec{
			Backend: nil,
			TLS:     nil,
			Rules: []networking.IngressRule{
				{
					testHost,
					networking.IngressRuleValue{
						HTTP: &networking.HTTPIngressRuleValue{
							Paths: []networking.HTTPIngressPath{
								{
									Path: expectedService.Paths[0],
									Backend: networking.IngressBackend{
										ServiceName: testServiceName,
										ServicePort: intstr.IntOrString{IntVal: int32(expectedService.Port)},
									},
								},
								{
									Path: expectedService.Paths[1],
									Backend: networking.IngressBackend{
										ServiceName: testServiceName,
										ServicePort: intstr.IntOrString{IntVal: int32(expectedService.Port)},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	ingress.Namespace = testNamespace
	actualService, err := translation.IngressToService(&ingress)

	if err != nil {
		t.Fatalf("error translating networking.Ingress to k8s.Service: %v", err)
	}

	if !servicesMatch(expectedService, *actualService) {
		t.Fatalf("expected Service to be: %v, got: %v", expectedService, *actualService)
	}
}

// Having am empty list of addresses is the most likely case when Ingress resources are created
// (as it takes k8s some amount of time to requisition the addresses).
// A subsequent Update event will contain the Address. The Registration.Modify method should handle that.
func TestTranslation_IngressToService_NoAddressesPresent(t *testing.T) {
	translation := new(Translation)
	testHost := "test-host"
	testNamespace := "testing-namespace"
	testServiceName := "test-service"
	var addresses []string
	expectedService := k8s.Service{
		Addresses: addresses,
		Name:      translation.FormatServiceName(testNamespace, testServiceName),
		Paths:     []string{"/apple", "/banana"},
		Port:      80,
	}

	ingress := networking.Ingress{
		Spec: networking.IngressSpec{
			Backend: nil,
			TLS:     nil,
			Rules: []networking.IngressRule{
				{
					testHost,
					networking.IngressRuleValue{
						HTTP: &networking.HTTPIngressRuleValue{
							Paths: []networking.HTTPIngressPath{
								{
									Path: expectedService.Paths[0],
									Backend: networking.IngressBackend{
										ServiceName: testServiceName,
										ServicePort: intstr.IntOrString{IntVal: int32(expectedService.Port)},
									},
								},
								{
									Path: expectedService.Paths[1],
									Backend: networking.IngressBackend{
										ServiceName: testServiceName,
										ServicePort: intstr.IntOrString{IntVal: int32(expectedService.Port)},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	ingress.Namespace = testNamespace
	actualService, err := translation.IngressToService(&ingress)

	if err != nil {
		t.Fatalf("error translating networking.Ingress to k8s.Service: %v", err)
	}

	if !servicesMatch(expectedService, *actualService) {
		t.Fatalf("expected Service to be: %v, got: %v", expectedService, *actualService)
	}
}
