package kong

import (
	"context"
	"fmt"
	"github.com/Sirupsen/logrus"
	gokong "github.com/hbagdi/go-kong/kong"
	"strings"
)

type Registrar interface {
	Register(service ServiceDef) error
	Deregister(service ServiceDef) error
	Modify(prevService ServiceDef, newService ServiceDef) error
}

type Registration struct {
	Kong        Client
	context     context.Context
	listOptions gokong.ListOpt
}

func NewRegistration(kongClient Client) (Registration, error) {
	registration := new(Registration)
	registration.Kong = kongClient
	return *registration, nil
}

func (registration *Registration) Deregister(service ServiceDef) error {
	return nil
}

func (registration *Registration) Register(serviceDef ServiceDef) error {
	var gerr error

	unacceptableHttpStatus := func(err error) bool {
		if err == nil {
			return false
		}

		if strings.Contains(err.Error(), "409") {
			logrus.Warnf("Skipping an error: %v", err)
			return false
		} else {
			return true
		}
	}

	kongService, _ := buildService(serviceDef)
	_, err := registration.Kong.Services.Create(registration.context, &kongService)
	if unacceptableHttpStatus(err) {
		gerr = fmt.Errorf("Registration::Register failed to create the Service: %v", err)
	}

	for routeName, path := range serviceDef.RoutesMap {
		route, err := buildRoute(kongService, routeName, path, false)
		if err != nil {
			gerr = fmt.Errorf("Registration::Register failed to build Route for path '%s': %v. Previous errors: %v", path, err, gerr)
		}

		_, err = registration.Kong.Routes.Create(registration.context, &route)
		if unacceptableHttpStatus(err) {
			gerr = fmt.Errorf("Registration::Register failed to create Route for path '%s': %v. Previous errors: %v", path, err, gerr)
		}
	}

	upstream, _ := buildUpstream(serviceDef)
	_, err = registration.Kong.Upstreams.Create(registration.context, &upstream)
	if unacceptableHttpStatus(err) {
		gerr = fmt.Errorf("Registration::Register failed to create Upstream: %v. Previous errors: %v", err, gerr)
	}

	for _, targetAddress := range serviceDef.Targets {
		target, err := buildTarget(upstream, targetAddress)
		if err != nil {
			gerr = fmt.Errorf("Registration::Register failed to build Target for address '%s': %v. Previous errors: %v", targetAddress, err, gerr)
		}

		_, err = registration.Kong.Targets.Create(registration.context, &target)
		if unacceptableHttpStatus(err) {
			gerr = fmt.Errorf("Registration::Register failed to create Target for address '%s': %v. Previous errors: %v", targetAddress, err, gerr)
		}
	}

	return gerr
}

func (registration *Registration) Modify(prevService ServiceDef, newService ServiceDef) error {
	return nil
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

func buildService(serviceDef ServiceDef) (gokong.Service, error) {
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

func buildTarget(upstream gokong.Upstream, targetAddress string) (gokong.Target, error) {
	target := gokong.Target{
		CreatedAt: nil,
		ID:        nil,
		Target:    &targetAddress,
		Upstream:  &upstream,
		Weight:    nil,
		Tags:      nil,
	}
	return target, nil
}

func buildUpstream(serviceDef ServiceDef) (gokong.Upstream, error) {
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
