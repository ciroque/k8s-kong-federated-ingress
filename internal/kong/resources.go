package kong

import (
	"fmt"
)

type ServiceDef struct {
	Addresses []string
	Name      string
	Names     ResourceNames
	Namespace string
	Paths     []string
	Port      int
}

type ResourceNames struct {
	RouteName    string
	ServiceName  string
	UpstreamName string
}

func NewResourceNames(service ServiceDef) ResourceNames {
	return ResourceNames{
		RouteName:    fmt.Sprintf("%s.%s.route", service.Namespace, service.Name),
		ServiceName:  fmt.Sprintf("%s.%s.service", service.Namespace, service.Name),
		UpstreamName: fmt.Sprintf("%s.%s.upstream", service.Namespace, service.Name),
	}
}
