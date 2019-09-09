package kong

import (
	"context"
	"github.com/ciroque/k8s-kong-federated-ingress/internal/k8s"
	"github.com/hbagdi/go-kong/kong"
)

type Registrar interface {
	Register(service *k8s.Service) (*RegisteredResources, error)
	Deregister(service *k8s.Service) error
	Modify(prevService *k8s.Service, newService *k8s.Service) (*RegisteredResources, error)
}

type Registration struct {
	Kong        *kong.Client
	context     context.Context
	listOptions kong.ListOpt
}

func (registration *Registration) Deregister(service *k8s.Service) error {
	return nil
}

func (registration *Registration) Register(service *k8s.Service) (*RegisteredResources, error) {
	return nil, nil
}

func (registration *Registration) Modify(prevService *k8s.Service, newService *k8s.Service) (*RegisteredResources, error) {
	return nil, nil
}
