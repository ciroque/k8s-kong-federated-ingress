package kong

import (
	"context"
	"fmt"
	gokong "github.com/hbagdi/go-kong/kong"
)

type Registrar interface {
	Register(service KongServiceDef) error
	Deregister(service KongServiceDef) error
	Modify(prevService KongServiceDef, newService KongServiceDef) error
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

func (registration *Registration) Deregister(service KongServiceDef) error {
	return nil
}

func (registration *Registration) Register(serviceDef KongServiceDef) error {

	var gerr error

	kongService, _ := buildService(serviceDef)
	_, err := registration.Kong.Services.Create(registration.context, kongService)
	if err != nil {
		return fmt.Errorf("Registration::Register failed to create the ServicesMap: %v", err)
	}

	for _, path := range serviceDef.Routes {
		routeName := buildRouteName(serviceDef.ServiceName, path)
		route, err := buildRoute(kongService, routeName, path, false)
		if err != nil {
			gerr = fmt.Errorf("Registration::Register failed to build Route for path '%s': %v. Previous errors: %v", path, err, gerr)
		}

		_, err = registration.Kong.Routes.Create(registration.context, route)
		if err != nil {
			gerr = fmt.Errorf("Registration::Register failed to create Route for path '%s': %v. Previous errors: %v", path, err, gerr)
		}
	}

	upstream, _ := buildUpstream(serviceDef)
	_, err = registration.Kong.Upstreams.Create(registration.context, upstream)
	if err != nil {
		return fmt.Errorf("Registration::Register failed to create Upstream: %v. Previous errors: %v", err, gerr)
	}

	for _, targetAddress := range serviceDef.Targets {
		target, err := buildTarget("upstreamTBD", targetAddress)
		if err != nil {
			gerr = fmt.Errorf("Registration::Register failed to build Target for address '%s': %v. Previous errors: %v", targetAddress, err, gerr)
		}

		_, err = registration.Kong.Targets.Create(registration.context, target)
		if err != nil {
			gerr = fmt.Errorf("Registration::Register failed to create Target for address '%s': %v. Previous errors: %v", targetAddress, err, gerr)
		}
	}

	return gerr
}

func (registration *Registration) Modify(prevService KongServiceDef, newService KongServiceDef) error {
	return nil
}

// TODO
func buildRouteName(serviceName string, route string) string {
	return "foo"
}

func buildRoute(service gokong.Service, routeName string, path string, stripPath bool) (gokong.Route, error) {
	kongRoute := gokong.Route{
		CreatedAt:               nil,
		Hosts:                   nil,
		Headers:                 nil,
		ID:                      nil,
		Name:                    &routeName,
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

func buildService(serviceDef KongServiceDef) (gokong.Service, error) {
	kongService := gokong.Service{
		ClientCertificate: nil,
		ConnectTimeout:    nil,
		CreatedAt:         nil,
		Host:              &serviceDef.UpstreamName,
		ID:                nil,
		Name:              &serviceDef.ServiceName,
		Path:              nil,
		Port:              nil,
		Protocol:          nil,
		ReadTimeout:       nil,
		Retries:           nil,
		UpdatedAt:         nil,
		WriteTimeout:      nil,
		Tags:              nil,
	}

	return kongService, nil
}

func buildTarget(upstream string, targetAddress string) (gokong.Target, error) {
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
func buildUpstream(serviceDef KongServiceDef) (gokong.Upstream, error) {
	upstream := gokong.Upstream{
		ID:                 nil,
		Name:               &serviceDef.UpstreamName,
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
