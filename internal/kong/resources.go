package kong

import (
	"fmt"
)

type ServiceDef struct {
	Addresses []string
	Name      string
	Paths     []string
	Port      int
}

type ResourceNames struct {
	ServiceName  string
	UpstreamName string
}

func NewResourceNames(service ServiceDef) ResourceNames {
	return ResourceNames{
		ServiceName:  fmt.Sprintf("%s.service", service.Name),
		UpstreamName: fmt.Sprintf("%s.upstream", service.Name),
	}
}
