package k8s

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
)

type KongServiceDef struct {
	Service  string
	Routes   []string
	Upstream string
	Targets  []string
}

type Translator interface {
	IngressToK8sService(ingress *networking.Ingress) (ServiceDef, error)
}

type Translation struct {
}

func (translation *Translation) IngressToK8sService(ingress *networking.Ingress) (ServiceDef, error) {
	serviceDef := new(ServiceDef)
	var paths []string

	for _, rule := range ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			serviceDef.Name = path.Backend.ServiceName // This is a problem if there are multiple paths pointing to different services...
			serviceDef.Port = int(path.Backend.ServicePort.IntVal)
			paths = append(paths, path.Path)
		}
	}

	serviceDef.Namespace = ingress.Namespace
	serviceDef.Paths = paths
	serviceDef.Addresses = buildAddresses(ingress.Status.LoadBalancer.Ingress)

	return *serviceDef, nil
}

func buildAddresses(ingresses []v1.LoadBalancerIngress) []string {
	var addresses []string
	for _, ingress := range ingresses {
		address := fmt.Sprintf("%s:80", ingress.IP)
		addresses = append(addresses, address)
	}
	return addresses
}
