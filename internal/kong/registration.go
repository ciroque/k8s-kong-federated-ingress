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

func (registration *Registration) Register(service ServiceDef) error {
	/// ServiceDef
	svc, _ := buildService(service)
	registration.Kong.Services.Create(registration.context, svc)

	/// Route

	/// Upstream

	upstream, _ := buildUpstream(service)
	_, err := registration.Kong.Upstreams.Create(registration.context, upstream)

	/// Targets

	return err
}

func (registration *Registration) Modify(prevService ServiceDef, newService ServiceDef) error {
	return nil
}

func buildService(service ServiceDef) (gokong.Service, error) {
	name := fmt.Sprintf("%s.upstream", service.Name)
	kongService := gokong.Service{
		ClientCertificate: nil,
		ConnectTimeout:    nil,
		CreatedAt:         nil,
		Host:              nil,
		ID:                nil,
		Name:              &name,
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

/// TODO: Support for Health Checks
func buildUpstream(service ServiceDef) (gokong.Upstream, error) {
	name := fmt.Sprintf("%s.upstream", service.Name)
	upstream := gokong.Upstream{
		ID:                 nil,
		Name:               &name,
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
