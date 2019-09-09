package kong

import (
	"github.com/ciroque/k8s-kong-federated-ingress/internal/k8s"
	networking "k8s.io/api/networking/v1beta1"
)

type Translator interface {
	IngressToService(ingress *networking.Ingress) (*k8s.Service, error)
}

type Translation struct {
}

func (translation *Translation) IngressToService(ingress *networking.Ingress) (*k8s.Service, error) {
	return nil, nil
}
