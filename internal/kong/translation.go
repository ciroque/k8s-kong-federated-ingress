package kong

import (
	"github.marchex.com/marchex/k8s-kong-federated-ingress/internal/k8s"
	"regexp"
)

type Translator interface {
	ServiceToKong(serviceName string, service k8s.ServiceDef) (ServiceDef, error)
}

type Translation struct {
}

func (translation *Translation) ServiceToKong(serviceName string, service k8s.ServiceDef) (ServiceDef, error) {
	kongServiceDef := ServiceDef{
		ServiceName:  service.Namespace + "." + serviceName + ".service",
		RoutesMap:    translation.buildRoutesMap(service.Namespace, serviceName, service.Paths),
		UpstreamName: service.Namespace + "." + serviceName + ".upstream",
		Targets:      service.Addresses,
	}
	return kongServiceDef, nil
}

func (translation Translation) buildRoutesMap(namespace string, serviceName string, paths []string) map[string]string {

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	var cleanFunction func(string) string
	if err == nil {
		cleanFunction = func(text string) string {
			return reg.ReplaceAllString(text, "")
		}
	} else {
		cleanFunction = func(text string) string {
			return text
		}
	}

	routeMap := make(map[string]string)
	for _, path := range paths {
		key := namespace + "." + serviceName + "." + cleanFunction(path) + ".route"
		routeMap[key] = path
	}
	return routeMap
}
