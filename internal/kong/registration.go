package kong

import (
	"context"
	"fmt"
	"github.com/ciroque/k8s-kong-federated-ingress/internal/k8s"
	gokong "github.com/hbagdi/go-kong/kong"
)

type Registrar interface {
	Register(service k8s.Service) error
	Deregister(service k8s.Service) error
	Modify(prevService k8s.Service, newService k8s.Service) error
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

func (registration *Registration) Deregister(service k8s.Service) error {
	return nil
}

func (registration *Registration) Register(service k8s.Service) error {
	/// Service
	svc, _ := buildService(service)
	registration.Kong.Services.Create(registration.context, svc)

	/// Route

	/// Upstream

	upstream, _ := buildUpstream(service)
	_, err := registration.Kong.Upstreams.Create(registration.context, upstream)

	/// Targets

	return err
}

func (registration *Registration) Modify(prevService k8s.Service, newService k8s.Service) error {
	return nil
}

func buildService(service k8s.Service) (gokong.Service, error) {
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
func buildUpstream(service k8s.Service) (gokong.Upstream, error) {
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
