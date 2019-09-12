package kong

import (
	"fmt"
	"github.com/ciroque/k8s-kong-federated-ingress/internal/k8s"
)

type ServiceDef struct {
	Service  string
	Routes   []string
	Upstream string
	Targets  []string
}

/// TODO DEPRECATE THIS
type ResourceNames struct {
	RouteName    string
	ServiceName  string
	UpstreamName string
}

func NewResourceNames(service k8s.ServiceDef) ResourceNames {
	return ResourceNames{
		RouteName:    fmt.Sprintf("%s.%s.route", service.Namespace, service.Name),
		ServiceName:  fmt.Sprintf("%s.%s.service", service.Namespace, service.Name),
		UpstreamName: fmt.Sprintf("%s.%s.upstream", service.Namespace, service.Name),
	}
}
