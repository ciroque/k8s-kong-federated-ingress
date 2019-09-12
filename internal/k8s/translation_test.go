package k8s

import (
	v1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

func TestTranslation_IngressToService(t *testing.T) {
	translation := new(Translation)
	testHost := "test-host"
	testNamespace := "testing-namespace"
	testServiceName := "test-service"
	addresses := []string{
		"10.200.30.400:80",
		"10.200.30.401:80",
	}

	expectedServiceMap := ServicesMap{
		testServiceName: ServiceDef{
			Addresses: addresses,
			Name:      testServiceName,
			Namespace: testNamespace,
			Paths:     []string{"/apple", "/banana"},
			Port:      80,
		},
	}

	expectedService := ServiceDef{
		Addresses: addresses,
		Name:      testServiceName,
		Namespace: testNamespace,
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
		t.Fatalf("error translating networking.Ingress to k8s.K8sServiceDef: %v", err)
	}

	if !ServicesMapMatch(expectedServiceMap, actualService) {
		t.Fatalf("expected K8sServiceDef to be:\n\t%v\n, got:\n\t%v", expectedServiceMap, actualService)
	}
}

func TestTranslation_IngressToServiceDef_MultipleServicesInRules(t *testing.T) {
	translation := new(Translation)
	testHost := "test-host"
	testNamespace := "testing-namespace"
	testServiceName1 := "test-service-1"
	testServiceName2 := "test-service-2"
	addresses := []string{
		"10.200.30.400:80",
		"10.200.30.401:80",
	}

	expectedServiceMap := ServicesMap{
		testServiceName1: ServiceDef{
			Addresses: addresses,
			Name:      testServiceName1,
			Namespace: testNamespace,
			Paths:     []string{"/apple", "/banana"},
			Port:      80,
		},
	}

	expectedService := ServiceDef{
		Addresses: addresses,
		Name:      testServiceName1,
		Namespace: testNamespace,
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
										ServiceName: testServiceName1,
										ServicePort: intstr.IntOrString{IntVal: int32(expectedService.Port)},
									},
								},
								{
									Path: expectedService.Paths[1],
									Backend: networking.IngressBackend{
										ServiceName: testServiceName2,
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
		t.Fatalf("error translating networking.Ingress to k8s.K8sServiceDef: %v", err)
	}

	if !ServicesMapMatch(expectedServiceMap, actualService) {
		t.Fatalf("expected K8sServiceDef to be:\n\t%v\n, got:\n\t%v", expectedService, actualService)
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
	expectedService := ServiceDef{
		Addresses: addresses,
		Name:      testServiceName,
		Namespace: testNamespace,
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
		t.Fatalf("error translating networking.Ingress to k8s.K8sServiceDef: %v", err)
	}

	if !ServicesMatch(expectedService, actualService) {
		t.Fatalf("expected K8sServiceDef to be:\n\t%v,\ngot:\n\t%v", expectedService, actualService)
	}
}
