package kong

import (
	"github.com/ciroque/k8s-kong-federated-ingress/internal/k8s"
)

type Translator interface {
	ServiceToKong(serviceName string, service k8s.ServiceDef) (KongServiceDef, error)
}

type Translation struct {
}

func (translation *Translation) ServiceToKong(serviceName string, service k8s.ServiceDef) (KongServiceDef, error) {
	kongServiceDef := KongServiceDef{
		ServiceName:  service.Namespace + "-" + serviceName + ".service",
		Routes:       service.Paths,
		UpstreamName: service.Namespace + "-" + serviceName + ".upstream",
		Targets:      service.Addresses,
	}
	return kongServiceDef, nil
}
