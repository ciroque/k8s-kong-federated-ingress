package kong

import (
	"context"
	"fmt"
	gokong "github.com/hbagdi/go-kong/kong"
)

type Registrar interface {
	Register(service ServiceDef) error
	Deregister(service ServiceDef) error
	Modify(prevService ServiceDef, newService ServiceDef) error
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

func (registration *Registration) Deregister(service ServiceDef) error {
	return nil
}

func (registration *Registration) Register(serviceDef ServiceDef) error {
	resourceNames := NewResourceNames(serviceDef)

	var gerr error

	/// Service
	kongService, _ := buildService(serviceDef, resourceNames)
	_, err := registration.Kong.Services.Create(registration.context, kongService)
	if err != nil {
		return fmt.Errorf("Registration::Register failed to create the Service: %v", err)
	}

	/// Routes
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

	return gerr
}

func (registration *Registration) Modify(prevService ServiceDef, newService ServiceDef) error {
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

func buildService(serviceDef ServiceDef, resourceNames ResourceNames) (gokong.Service, error) {
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

/// TODO: Support for Health Checks
func buildUpstream(_ ServiceDef, resourceNames ResourceNames) (gokong.Upstream, error) {
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
