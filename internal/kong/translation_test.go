package kong

import (
	"github.com/ciroque/k8s-kong-federated-ingress/internal/k8s"
	networking "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

func compareStringArrays(l []string, r []string) bool {
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
		l.Port == r.Port /// && compareStringArrays(l.Addresses, r.Addresses)
}

func TestTranslation_IngressToService_RouteAndService(t *testing.T) {
	expectedService := k8s.Service{
		Addresses: nil,
		Name:      "test-service-name",
		Paths:     []string{"/apple", "/banana"},
		Port:      80,
	}

	translation := new(Translation)
	ingress := networking.Ingress{
		Spec: networking.IngressSpec{
			Rules: []networking.IngressRule{
				{
					"test-host",
					networking.IngressRuleValue{
						HTTP: &networking.HTTPIngressRuleValue{
							Paths: []networking.HTTPIngressPath{
								{
									Path: expectedService.Paths[0],
									Backend: networking.IngressBackend{
										ServiceName: expectedService.Name,
										ServicePort: intstr.IntOrString{IntVal: int32(expectedService.Port)},
									},
								},
								{
									Path: expectedService.Paths[1],
									Backend: networking.IngressBackend{
										ServiceName: expectedService.Name,
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
	actualService, err := translation.IngressToService(&ingress)

	if err != nil {
		t.Fatalf("error translating networking.Ingress to k8s.Service: %v", err)
	}

	if !servicesMatch(expectedService, *actualService) {
		t.Fatalf("expected Service to be: %v, got: %v", expectedService, *actualService)
	}
}
