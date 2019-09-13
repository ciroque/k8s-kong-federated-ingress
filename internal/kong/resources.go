/*

	DEPRECATING THIS ; THE NEW HOTNESS IS Translation / Registration implementation

*/

package kong

import (
	"fmt"
	"github.com/ciroque/k8s-kong-federated-ingress/internal/k8s"
)

/// TODO DEPRECATE THIS
type ResourceNames struct {
	RouteName    string
	ServiceName  string
	UpstreamName string
}

func NewResourceNames(service k8s.ServiceDef) ResourceNames {
	return ResourceNames{
		RouteName:    fmt.Sprintf("%s.%s.route", service.Namespace, "name"),
		ServiceName:  fmt.Sprintf("%s.%s.service", service.Namespace, "name"),
		UpstreamName: fmt.Sprintf("%s.%s.upstream", service.Namespace, "name"),
	}
}
