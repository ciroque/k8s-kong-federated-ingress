package kong

import (
	"context"
	"errors"
	"fmt"
	gokong "github.com/hbagdi/go-kong/kong"
	"reflect"
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

type TestRoutes struct {
	Created *[]gokong.Route
}

func (routes TestRoutes) Create(context context.Context, route *gokong.Route) (*gokong.Route, error) {
	*routes.Created = append(*routes.Created, *route)
	return route, nil
}

type FailRoutes struct {
}

func (routes FailRoutes) Create(context context.Context, route *gokong.Route) (*gokong.Route, error) {
	return route, errors.New("420 Enhance your calm")
}

type TestServices struct {
	Service *gokong.Service
}

func (services TestServices) Create(context context.Context, service *gokong.Service) (*gokong.Service, error) {
	*services.Service = *service
	return service, nil
}

type FailServices struct {
}

func (streams FailServices) Create(context context.Context, service *gokong.Service) (*gokong.Service, error) {
	return service, errors.New("420 Enhance your calm")
}

type TestTargets struct {
	Created *[]gokong.Target
}

func (targets TestTargets) Create(context context.Context, target *gokong.Target) (*gokong.Target, error) {
	*targets.Created = append(*targets.Created, *target)
	return target, nil
}

type FailTargets struct {
}

func (targets FailTargets) Create(context context.Context, target *gokong.Target) (*gokong.Target, error) {
	return target, errors.New("420 Enhance your calm")
}

type TestUpstreams struct {
	Created *gokong.Upstream
}

func (upstreams TestUpstreams) Create(context context.Context, upstream *gokong.Upstream) (*gokong.Upstream, error) {
	*upstreams.Created = *upstream
	return upstream, nil
}

type FailUpstreams struct {
}

func (upstreams FailUpstreams) Create(context context.Context, upstream *gokong.Upstream) (*gokong.Upstream, error) {
	return upstream, errors.New("420 Enhance your calm")
}

/// ********************************************************************************************************************
/// TESTS

func TestRegistration_Register_CreateRouteCalled(t *testing.T) {
	mockClient, routes, service, _, _ := buildMockClient()

	registration, _ := NewRegistration(Client(mockClient))
	serviceDef := buildExampleServiceDef()

	err := registration.Register(serviceDef)
	if err != nil {
		t.Fatalf("Register failed with: %v", err)
	}

	routeNames := []string{}
	routePaths := []string{}
	for name, path := range serviceDef.RoutesMap {
		routeNames = append(routeNames, name)
		routePaths = append(routePaths, path)
	}

	stripPath := false

	expectedRoutes := []gokong.Route{
		{
			Name:      &routeNames[0],
			Paths:     gokong.StringSlice(routePaths[0]),
			Service:   service.Service,
			StripPath: &stripPath,
		},
		{
			Name:      &routeNames[1],
			Paths:     gokong.StringSlice(routePaths[1]),
			Service:   service.Service,
			StripPath: &stripPath,
		},
	}

	if !reflect.DeepEqual(expectedRoutes, *routes.Created) {
		t.Fatal(fmt.Sprintf("Expected TestRoutes.Create to be called with:\n\t%v, \nbut got\n\t%v", expectedRoutes, *routes.Created))
	}
}

func TestRegistration_Register_CreateRouteFails(t *testing.T) {
	mockClient, _, _, _, _ := buildMockClient()
	mockClient.Routes = FailRoutes{}

	registration, _ := NewRegistration(Client(mockClient))
	serviceDef := buildExampleServiceDef()

	err := registration.Register(serviceDef)
	if err == nil {
		t.Fatalf("Register should have failed.")
	}

	if !strings.Contains(err.Error(), "420") {
		t.Fatalf("Failure message should contain a '420 Enhance your calm' message. Got: %v", err)
	}
}

func TestRegistration_Register_CreateServiceCalled(t *testing.T) {
	mockClient, _, services, _, _ := buildMockClient()

	registration, _ := NewRegistration(Client(mockClient))
	serviceDef := buildExampleServiceDef()
	expectedService := gokong.Service{
		Host: &serviceDef.UpstreamName,
		Name: &serviceDef.ServiceName,
	}

	err := registration.Register(serviceDef)
	if err != nil {
		t.Fatalf("Register failed with: %v", err)
	}

	if !reflect.DeepEqual(expectedService, *services.Service) {
		t.Fatal(fmt.Sprintf("Expected TestServices.Create to be called with:\n\t%v, \nbut got\n\t%v", expectedService, *services.Service))
	}
}

func TestRegistration_Register_CreateServiceFails(t *testing.T) {
	mockClient, _, _, _, _ := buildMockClient()
	mockClient.Services = FailServices{}

	registration, _ := NewRegistration(Client(mockClient))
	serviceDef := buildExampleServiceDef()

	err := registration.Register(serviceDef)
	if err == nil {
		t.Fatalf("Register should have failed.")
	}

	if !strings.Contains(err.Error(), "420") {
		t.Fatalf("Failure message should contain a '420 Enhance your calm' message. Got: %v", err)
	}
}

func TestRegistration_Register_CreateTargetCalled(t *testing.T) {
	mockClient, _, _, targets, upstreams := buildMockClient()

	registration, _ := NewRegistration(Client(mockClient))
	serviceDef := buildExampleServiceDef()
	expectedTargets := []gokong.Target{
		{
			CreatedAt: nil,
			ID:        nil,
			Target:    &serviceDef.Targets[0],
			Upstream:  upstreams.Created,
			Weight:    nil,
			Tags:      nil,
		},
		{
			CreatedAt: nil,
			ID:        nil,
			Target:    &serviceDef.Targets[1],
			Upstream:  upstreams.Created,
			Weight:    nil,
			Tags:      nil,
		},
		{
			CreatedAt: nil,
			ID:        nil,
			Target:    &serviceDef.Targets[2],
			Upstream:  upstreams.Created,
			Weight:    nil,
			Tags:      nil,
		},
	}

	err := registration.Register(serviceDef)
	if err != nil {
		t.Fatalf("Register failed with: %v", err)
	}

	if !reflect.DeepEqual(expectedTargets, *targets.Created) {
		t.Fatal(fmt.Sprintf("Expected TestTargets.Create to be called with:\n\t%v, \nbut got\n\t%v", expectedTargets, *targets.Created))
	}
}

func TestRegistration_Register_CreateTargetFails(t *testing.T) {
	mockClient, _, _, _, _ := buildMockClient()
	mockClient.Targets = FailTargets{}

	registration, _ := NewRegistration(Client(mockClient))
	serviceDef := buildExampleServiceDef()

	err := registration.Register(serviceDef)
	if err == nil {
		t.Fatalf("Register should have failed.")
	}

	if !strings.Contains(err.Error(), "420") {
		t.Fatalf("Failure message should contain a '420 Enhance your calm' message. Got: %v", err)
	}
}

func TestRegistration_Register_CreateUpstreamCalled(t *testing.T) {
	mockClient, _, _, _, upstreams := buildMockClient()

	registration, _ := NewRegistration(Client(mockClient))
	serviceDef := buildExampleServiceDef()
	expectedUpstream := gokong.Upstream{
		Name: &serviceDef.UpstreamName,
	}

	err := registration.Register(serviceDef)
	if err != nil {
		t.Fatalf("Register failed with: %v", err)
	}

	if !reflect.DeepEqual(expectedUpstream, *upstreams.Created) {
		t.Fatal(fmt.Sprintf("Expected TestUpstreams.Create to be called with:\n\t%v, \nbut got\n\t%v", expectedUpstream, *upstreams.Created))
	}
}

func TestRegistration_Register_CreateUpstreamFails(t *testing.T) {
	mockClient, _, _, _, _ := buildMockClient()
	mockClient.Upstreams = FailUpstreams{}

	registration, _ := NewRegistration(Client(mockClient))
	serviceDef := buildExampleServiceDef()

	err := registration.Register(serviceDef)
	if err == nil {
		t.Fatalf("Register should have failed.")
	}

	if !strings.Contains(err.Error(), "420") {
		t.Fatalf("Failure message should contain a '420 Enhance your calm' message. Got: %v", err)
	}
}

/// ********************************************************************************************************************
/// HELPERS

func buildMockClient() (Client, TestRoutes, TestServices, TestTargets, TestUpstreams) {
	emptyRoutes := []gokong.Route{}
	emptyTargets := []gokong.Target{}

	routes := TestRoutes{Created: &emptyRoutes}
	services := TestServices{
		Service: new(gokong.Service),
	}
	targets := TestTargets{Created: &emptyTargets}
	upstreams := TestUpstreams{Created: new(gokong.Upstream)}

	mockClient := MockClient{
		Routes:    routes,
		Services:  services,
		Targets:   targets,
		Upstreams: upstreams,
	}
	return Client(mockClient), routes, services, targets, upstreams
}

func buildExampleServiceDef() ServiceDef {
	serviceDef := ServiceDef{
		ServiceName: "test-service",
		RoutesMap: map[string]string{
			"test-service.10.100.100.10.route": "10.100.100.10",
			"test-service.10.100.100.11.route": "10.100.100.11",
		},
		UpstreamName: "test-service.upstream",
		Targets: []string{
			"/apples",
			"/bananas",
			"/oranges",
		},
	}
	return serviceDef
}
