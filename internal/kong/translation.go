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

	k8sService := new(k8s.Service)
	var paths []string

	for _, rule := range ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			k8sService.Name = path.Backend.ServiceName
			k8sService.Port = int(path.Backend.ServicePort.IntVal)
			paths = append(paths, path.Path)
		}
	}

	k8sService.Paths = paths

	return k8sService, nil
}
