package kong

import (
	"fmt"
	"github.com/ciroque/k8s-kong-federated-ingress/internal/k8s"
	v1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
)

type Translator interface {
	IngressToService(ingress *networking.Ingress) (*k8s.Service, error)
}

type Translation struct {
}

func (translation *Translation) FormatServiceName(namespace string, service string) string {
	return fmt.Sprintf("%s.%s.service", namespace, service)
}

func (translation *Translation) IngressToService(ingress *networking.Ingress) (*k8s.Service, error) {
	k8sService := new(k8s.Service)
	var paths []string

	for _, rule := range ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			k8sService.Name = translation.FormatServiceName(ingress.Namespace, path.Backend.ServiceName)
			k8sService.Port = int(path.Backend.ServicePort.IntVal)
			paths = append(paths, path.Path)
		}
	}

	k8sService.Paths = paths
	k8sService.Addresses = buildAddresses(ingress.Status.LoadBalancer.Ingress)

	return k8sService, nil
}

func buildAddresses(ingresses []v1.LoadBalancerIngress) []string {
	var addresses []string
	for _, ingress := range ingresses {
		address := fmt.Sprintf("%s:80", ingress.IP)
		addresses = append(addresses, address)
	}
	return addresses
}
