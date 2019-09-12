package kong

import (
	"github.com/ciroque/k8s-kong-federated-ingress/internal/k8s"
)

type Translator interface {
	ServiceToKong(service k8s.ServicesMap) (ServiceDef, error)
}
