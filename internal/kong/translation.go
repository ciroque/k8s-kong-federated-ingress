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
	return KongServiceDef{}, nil
}
