package k8s

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
)

type Translator interface {
	IngressToService(ingress *networking.Ingress) (ServiceDef, error)
}

type Translation struct {
}

func (translation *Translation) IngressToService(ingress *networking.Ingress) (ServicesMap, error) {
	servicesMap := make(map[string]ServiceDef)

	ingressAddresses := buildAddresses(ingress.Status.LoadBalancer.Ingress)

	for _, rule := range ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			serviceName := path.Backend.ServiceName
			serviceDef, found := servicesMap[serviceName]
			if !found {
				servicesMap[serviceName] = ServiceDef{
					Addresses: ingressAddresses,
					Namespace: ingress.Namespace,
					Paths:     []string{path.Path},
					Port:      int(path.Backend.ServicePort.IntVal),
				}
			} else {
				serviceDef.Paths = append(serviceDef.Paths, path.Path)
				servicesMap[serviceName] = serviceDef
			}
		}
	}

	return servicesMap, nil
}

func buildAddresses(ingresses []v1.LoadBalancerIngress) []string {
	var addresses []string
	for _, ingress := range ingresses {
		address := fmt.Sprintf("%s:80", ingress.IP)
		addresses = append(addresses, address)
	}
	return addresses
}
