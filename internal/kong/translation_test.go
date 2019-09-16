package kong

import (
	"github.com/ciroque/k8s-kong-federated-ingress/internal/k8s"
	"reflect"
	"testing"
)

func TestTranslation_K8sServiceToKongService(t *testing.T) {
	translation := new(Translation)

	serviceName := "aservice"

	serviceDef := k8s.ServiceDef{
		Addresses: []string{
			"10.200.30.400:80",
			"10.200.30.401:80",
		},
		Namespace: "namespace",
		Paths:     []string{"/apple", "/banana"},
	}

	expectedKongServiceDef := ServiceDef{
		ServiceName: serviceDef.Namespace + "-" + serviceName + ".service",
		RoutesMap: map[string]string{
			serviceDef.Namespace + "-" + serviceName + "-" + serviceDef.Paths[0] + ".route": serviceDef.Paths[0],
			serviceDef.Namespace + "-" + serviceName + "-" + serviceDef.Paths[1] + ".route": serviceDef.Paths[1],
		},
		UpstreamName: serviceDef.Namespace + "-" + serviceName + ".upstream",
		Targets:      serviceDef.Addresses,
	}

	actualKongServiceDef, err := translation.ServiceToKong(serviceName, serviceDef)

	if err != nil {
		t.Fatalf("error translating k8s.K8sServiceDef to kong.ServiceDef: %v", err)
	}

	if !reflect.DeepEqual(expectedKongServiceDef, actualKongServiceDef) {
		t.Fatalf("expected K8sServiceDef to be:\n\t%v\n, got:\n\t%v", expectedKongServiceDef, actualKongServiceDef)
	}
}
