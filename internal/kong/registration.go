package kong

import (
	"context"
	"fmt"
	"github.com/ciroque/k8s-kong-federated-ingress/internal/k8s"
	gokong "github.com/hbagdi/go-kong/kong"
)

type Registrar interface {
	Register(service k8s.ServiceDef) error
	Deregister(service k8s.ServiceDef) error
	Modify(prevService k8s.ServiceDef, newService k8s.ServiceDef) error
}

type Registration struct {
	Kong        ClientInterface
	context     context.Context
	listOptions gokong.ListOpt
}

func NewRegistration(kongClient ClientInterface) (Registration, error) {
	registration := new(Registration)
	registration.Kong = kongClient
	return *registration, nil
}

func (registration *Registration) Deregister(service k8s.ServiceDef) error {
	return nil
}

func (registration *Registration) Register(serviceDef k8s.ServiceDef) error {
	resourceNames := NewResourceNames(serviceDef)

	var gerr error

	/// Service
	kongService, _ := buildService(serviceDef, resourceNames)
	_, err := registration.Kong.Services.Create(registration.context, kongService)
	if err != nil {
		return fmt.Errorf("Registration::Register failed to create the Service: %v", err)
	}

	/// Routes
	/// TODO Route names need to be uniqueified
	for _, path := range serviceDef.Paths {
		route, err := buildRoute(kongService, resourceNames, path, false)
		if err != nil {
			gerr = fmt.Errorf("Registration::Register failed to build Route for path '%s': %v. Previous errors: %v", path, err, gerr)
		}

		_, err = registration.Kong.Routes.Create(registration.context, route)
		if err != nil {
			gerr = fmt.Errorf("Registration::Register failed to create Route for path '%s': %v. Previous errors: %v", path, err, gerr)
		}
	}

	/// Upstream
	upstream, _ := buildUpstream(serviceDef, resourceNames)
	_, err = registration.Kong.Upstreams.Create(registration.context, upstream)
	if err != nil {
		return fmt.Errorf("Registration::Register failed to create Upstream: %v. Previous errors: %v", err, gerr)
	}

	/// Targets
	//serviceDef.Addresses become Targets
	for _, address := range serviceDef.Addresses {
		target, err := buildTarget(serviceDef, resourceNames, address)
		if err != nil {
			gerr = fmt.Errorf("Registration::Register failed to build Target for address '%s': %v. Previous errors: %v", address, err, gerr)
		}

		_, err = registration.Kong.Targets.Create(registration.context, target)
		if err != nil {
			gerr = fmt.Errorf("Registration::Register failed to create Target for address '%s': %v. Previous errors: %v", address, err, gerr)
		}
	}

	return gerr
}

func (registration *Registration) Modify(prevService k8s.ServiceDef, newService k8s.ServiceDef) error {
	return nil
}

func buildRoute(service gokong.Service, resourceNames ResourceNames, path string, stripPath bool) (gokong.Route, error) {
	kongRoute := gokong.Route{
		CreatedAt:               nil,
		Hosts:                   nil,
		Headers:                 nil,
		ID:                      nil,
		Name:                    &resourceNames.RouteName,
		Methods:                 nil,
		Paths:                   gokong.StringSlice(path),
		PreserveHost:            nil,
		Protocols:               nil,
		RegexPriority:           nil,
		Service:                 &service,
		StripPath:               &stripPath,
		UpdatedAt:               nil,
		SNIs:                    nil,
		Sources:                 nil,
		Destinations:            nil,
		Tags:                    nil,
		HTTPSRedirectStatusCode: nil,
	}
	return kongRoute, nil
}

func buildService(serviceDef k8s.ServiceDef, resourceNames ResourceNames) (gokong.Service, error) {
	kongService := gokong.Service{
		ClientCertificate: nil,
		ConnectTimeout:    nil,
		CreatedAt:         nil,
		Host:              &resourceNames.UpstreamName,
		ID:                nil,
		Name:              &resourceNames.ServiceName,
		Path:              nil,
		Port:              &serviceDef.Port,
		Protocol:          nil,
		ReadTimeout:       nil,
		Retries:           nil,
		UpdatedAt:         nil,
		WriteTimeout:      nil,
		Tags:              nil,
	}

	return kongService, nil
}

func buildTarget(serviceDef k8s.ServiceDef, names ResourceNames, s string) (gokong.Target, error) {
	target := gokong.Target{
		CreatedAt: nil,
		ID:        nil,
		Target:    nil,
		Upstream:  nil,
		Weight:    nil,
		Tags:      nil,
	}
	return target, nil
}

/// TODO: Support for Health Checks
func buildUpstream(_ k8s.ServiceDef, resourceNames ResourceNames) (gokong.Upstream, error) {
	upstream := gokong.Upstream{
		ID:                 nil,
		Name:               &resourceNames.UpstreamName,
		Algorithm:          nil,
		Slots:              nil,
		Healthchecks:       nil,
		CreatedAt:          nil,
		HashOn:             nil,
		HashFallback:       nil,
		HashOnHeader:       nil,
		HashFallbackHeader: nil,
		HashOnCookie:       nil,
		HashOnCookiePath:   nil,
		Tags:               nil,
	}

	return upstream, nil
}
