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

	/// Service
	svc, _ := buildService(serviceDef, resourceNames)
	_, err := registration.Kong.Services.Create(registration.context, svc)
	if err != nil {
		return fmt.Errorf("Registration::Register failed to create the Service: %v", err)
	}

	/// Route

	/// Upstream
	upstream, _ := buildUpstream(serviceDef, resourceNames)
	_, err = registration.Kong.Upstreams.Create(registration.context, upstream)
	if err != nil {
		return fmt.Errorf("Registration::Register failed to create the Upstream: %v", err)
	}

	/// Targets

	return err
}

func (registration *Registration) Modify(prevService ServiceDef, newService ServiceDef) error {
	return nil
}

func buildService(serviceDef ServiceDef, names ResourceNames) (gokong.Service, error) {
	kongService := gokong.Service{
		ClientCertificate: nil,
		ConnectTimeout:    nil,
		CreatedAt:         nil,
		Host:              &names.UpstreamName,
		ID:                nil,
		Name:              &names.ServiceName,
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
func buildUpstream(_ ServiceDef, names ResourceNames) (gokong.Upstream, error) {
	upstream := gokong.Upstream{
		ID:                 nil,
		Name:               &names.UpstreamName,
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
