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
	ServiceName  string
	UpstreamName string
}

func NewResourceNames(service ServiceDef) ResourceNames {
	return ResourceNames{
		ServiceName:  fmt.Sprintf("%s.%s.service", service.Namespace, service.Name),
		UpstreamName: fmt.Sprintf("%s.%s.upstream", service.Namespace, service.Name),
	}
}
