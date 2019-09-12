package kong

import (
	"context"
	"errors"
	"fmt"
	gokong "github.com/hbagdi/go-kong/kong"
	"strings"
	"testing"
)

/// ********************************************************************************************************************
/// MOCKS

type MockClient struct {
	Routes    RoutesInterface
	Services  ServicesInterface
	Targets   TargetsInterface
	Upstreams UpstreamsInterface
}

type Routes struct {
	CreateCount *int
}

func (routes Routes) Create(context context.Context, route gokong.Route) (gokong.Route, error) {
	*routes.CreateCount = *routes.CreateCount + 1
	return route, nil
}

type FailRoutes struct {
	CreateCount *int
}

func (routes FailRoutes) Create(context context.Context, route gokong.Route) (gokong.Route, error) {
	*routes.CreateCount = *routes.CreateCount + 1
	return route, errors.New("409 Conflict")
}

type Services struct {
	CreateCount *int
}

func (streams Services) Create(context context.Context, service gokong.Service) (gokong.Service, error) {
	*streams.CreateCount = *streams.CreateCount + 1
	return service, nil
}

type FailServices struct {
	CreateCount *int
}

func (streams FailServices) Create(context context.Context, service gokong.Service) (gokong.Service, error) {
	*streams.CreateCount = *streams.CreateCount + 1
	return service, errors.New("409 Conflict")
}

type Targets struct {
	CreateCount *int
}

func (targets Targets) Create(context context.Context, target gokong.Target) (gokong.Target, error) {
	*targets.CreateCount = *targets.CreateCount + 1
	return target, nil
}

type FailTargets struct {
	CreateCount *int
}

func (targets FailTargets) Create(context context.Context, target gokong.Target) (gokong.Target, error) {
	*targets.CreateCount = *targets.CreateCount + 1
	return target, errors.New("409 Conflict")
}

type Upstreams struct {
	CreateCount *int
}

func (upstreams Upstreams) Create(context context.Context, upstream gokong.Upstream) (gokong.Upstream, error) {
	*upstreams.CreateCount = *upstreams.CreateCount + 1
	return upstream, nil
}

type FailUpstreams struct {
	CreateCount *int
}

func (upstreams FailUpstreams) Create(context context.Context, upstream gokong.Upstream) (gokong.Upstream, error) {
	*upstreams.CreateCount = *upstreams.CreateCount + 1
	return upstream, errors.New("409 Conflict")
}

/// ********************************************************************************************************************
/// TESTS

func TestRegistration_Register_CreateRouteCalled(t *testing.T) {
	mockClient, routes, _, _, _ := buildMockClient()

	registration, _ := NewRegistration(ClientInterface(mockClient))
	serviceDef := buildServiceDef()

	err := registration.Register(serviceDef)
	if err != nil {
		t.Fatalf("Register failed with: %v", err)
	}

	expectedCount := len(serviceDef.Paths)
	if *routes.CreateCount != expectedCount {
		t.Fatal(fmt.Sprintf("Expected Routes.Create to have been called %v times. Actual count: %v.", expectedCount, *routes.CreateCount))
	}
}

func TestRegistration_Register_CreateRouteFails(t *testing.T) {
	mockClient, _, _, _, _ := buildMockClient()
	mockClient.Routes = FailRoutes{CreateCount: new(int)} // Inject a mock that will fail the call

	registration, _ := NewRegistration(ClientInterface(mockClient))
	serviceDef := buildServiceDef()

	err := registration.Register(serviceDef)
	if err == nil {
		t.Fatalf("Register should have failed.")
	}

	if !strings.Contains(err.Error(), "409") {
		t.Fatalf("Failure message should contain a '409 Conflict' message. Got: %v", err)
	}
}

func TestRegistration_Register_CreateServiceCalled(t *testing.T) {
	mockClient, _, services, _, _ := buildMockClient()

	registration, _ := NewRegistration(ClientInterface(mockClient))
	serviceDef := buildServiceDef()

	err := registration.Register(serviceDef)
	if err != nil {
		t.Fatalf("Register failed with: %v", err)
	}

	if *services.CreateCount != 1 {
		t.Fatal("Expected Services.Create to have been called once. Actual count: ", *services.CreateCount)
	}
}

func TestRegistration_Register_CreateServiceFails(t *testing.T) {
	mockClient, _, _, _, _ := buildMockClient()
	mockClient.Services = FailServices{CreateCount: new(int)}

	registration, _ := NewRegistration(ClientInterface(mockClient))
	serviceDef := buildServiceDef()

	err := registration.Register(serviceDef)
	if err == nil {
		t.Fatalf("Register should have failed.")
	}

	if !strings.Contains(err.Error(), "409") {
		t.Fatalf("Failure message should contain a '409 Conflict' message. Got: %v", err)
	}
}

func TestRegistration_Register_CreateTargetCalled(t *testing.T) {
	mockClient, _, _, targets, _ := buildMockClient()

	registration, _ := NewRegistration(ClientInterface(mockClient))
	serviceDef := buildServiceDef()

	err := registration.Register(serviceDef)
	if err != nil {
		t.Fatalf("Register failed with: %v", err)
	}

	expectedCount := len(serviceDef.Addresses)
	if *targets.CreateCount != expectedCount {
		t.Fatal(fmt.Sprintf("Expected Routes.Create to have been called %v times. Actual count: %v.", expectedCount, *targets.CreateCount))
	}
}

func TestRegistration_Register_CreateTargetFails(t *testing.T) {
	mockClient, _, _, _, _ := buildMockClient()
	mockClient.Targets = FailTargets{CreateCount: new(int)} // Inject a mock that will fail the call

	registration, _ := NewRegistration(ClientInterface(mockClient))
	serviceDef := buildServiceDef()

	err := registration.Register(serviceDef)
	if err == nil {
		t.Fatalf("Register should have failed.")
	}

	if !strings.Contains(err.Error(), "409") {
		t.Fatalf("Failure message should contain a '409 Conflict' message. Got: %v", err)
	}
}

func TestRegistration_Register_CreateUpstreamCalled(t *testing.T) {
	mockClient, _, _, _, upstreams := buildMockClient()

	registration, _ := NewRegistration(ClientInterface(mockClient))
	serviceDef := buildServiceDef()

	err := registration.Register(serviceDef)
	if err != nil {
		t.Fatalf("Register failed with: %v", err)
	}

	if *upstreams.CreateCount != 1 {
		t.Fatal("Expected Upstreams.Create to have been called once. Actual count: ", *upstreams.CreateCount)
	}
}

func TestRegistration_Register_CreateUpstreamFails(t *testing.T) {
	mockClient, _, _, _, _ := buildMockClient()
	mockClient.Upstreams = FailUpstreams{CreateCount: new(int)}

	registration, _ := NewRegistration(ClientInterface(mockClient))
	serviceDef := buildServiceDef()

	err := registration.Register(serviceDef)
	if err == nil {
		t.Fatalf("Register should have failed.")
	}

	if !strings.Contains(err.Error(), "409") {
		t.Fatalf("Failure message should contain a '409 Conflict' message. Got: %v", err)
	}
}

/// ********************************************************************************************************************
/// HELPERS

func buildMockClient() (ClientInterface, Routes, Services, Targets, Upstreams) {
	routes := Routes{CreateCount: new(int)}
	services := Services{CreateCount: new(int)}
	targets := Targets{CreateCount: new(int)}
	upstreams := Upstreams{CreateCount: new(int)}

	mockClient := MockClient{
		Routes:    routes,
		Services:  services,
		Targets:   targets,
		Upstreams: upstreams,
	}
	return ClientInterface(mockClient), routes, services, targets, upstreams
}

func buildServiceDef() ServiceDef {
	serviceDef := ServiceDef{
		Addresses: []string{
			"10.100.100.10",
			"10.100.100.11",
			"10.100.100.12",
			"10.100.100.13",
			"10.100.100.14",
		},
		Name: "test-service",
		Paths: []string{
			"/apples",
			"/bananas",
			"/oranges",
		},
		Port: 8080,
	}
	serviceDef.Namespace = "test-namespace"
	return serviceDef
}
