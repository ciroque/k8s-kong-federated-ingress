package kong

import (
	"context"
	"fmt"
	gokong "github.com/hbagdi/go-kong/kong"
	"testing"
)

type MockClient struct {
	Routes    RoutesInterface
	Services  ServicesInterface
	Upstreams UpstreamsInterface
}

type Routes struct {
	CreateCount *int
}

func (routes Routes) Create(ctx context.Context, route gokong.Route) (gokong.Route, error) {
	*routes.CreateCount = *routes.CreateCount + 1
	return route, nil
}

type Services struct {
	CreateCount *int
}

func (streams Services) Create(ctx context.Context, service gokong.Service) (gokong.Service, error) {
	*streams.CreateCount = *streams.CreateCount + 1
	return service, nil
}

type Upstreams struct {
	CreateCount *int
}

func (upstreams Upstreams) Create(context context.Context, upstream gokong.Upstream) (gokong.Upstream, error) {
	*upstreams.CreateCount = *upstreams.CreateCount + 1
	return upstream, nil
}

func TestRegistration_Register_CreateUpstreamCalled(t *testing.T) {
	mockClient, _, _, upstreams := buildMockClient()

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

func TestRegistration_Register_CreateServiceCalled(t *testing.T) {
	mockClient, _, services, _ := buildMockClient()

	registration, _ := NewRegistration(ClientInterface(mockClient))
	serviceDef := buildServiceDef()

	err := registration.Register(serviceDef)
	if err != nil {
		t.Fatalf("Register failed with: %v", err)
	}

	if *services.CreateCount != 1 {
		t.Fatal("Expected Routes.Create to have been called once. Actual count: ", *services.CreateCount)
	}
}

func TestRegistration_Register_CreateRouteCalled(t *testing.T) {
	mockClient, routes, _, _ := buildMockClient()

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

func buildMockClient() (ClientInterface, Routes, Services, Upstreams) {
	routes := Routes{CreateCount: new(int)}
	services := Services{CreateCount: new(int)}
	upstreams := Upstreams{CreateCount: new(int)}

	mockClient := MockClient{
		Routes:    routes,
		Services:  services,
		Upstreams: upstreams,
	}
	return ClientInterface(mockClient), routes, services, upstreams
}

func buildServiceDef() ServiceDef {
	serviceDef := ServiceDef{
		Addresses: []string{},
		Name:      "test-service",
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
