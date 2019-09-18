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
	Deleted *[]string
}

func (routes TestRoutes) Create(context context.Context, route *gokong.Route) (*gokong.Route, error) {
	*routes.Created = append(*routes.Created, *route)
	return route, nil
}

func (routes TestRoutes) Delete(context context.Context, routeNameOrId *string) error {
	*routes.Deleted = append(*routes.Deleted, *routeNameOrId)
	return nil
}

type FailRoutes struct {
}

func (routes FailRoutes) Create(context context.Context, route *gokong.Route) (*gokong.Route, error) {
	return route, errors.New("420 Enhance your calm")
}

func (routes FailRoutes) Delete(context context.Context, routeNameOrId *string) error {
	return nil
}

type ConflictRoutes struct {
}

func (routes ConflictRoutes) Create(context context.Context, route *gokong.Route) (*gokong.Route, error) {
	return route, errors.New("409 Conflict")
}

func (routes ConflictRoutes) Delete(context context.Context, routeNameOrId *string) error {
	return nil
}

type TestServices struct {
	Service *gokong.Service
	Deleted *string
}

func (services TestServices) Create(context context.Context, service *gokong.Service) (*gokong.Service, error) {
	*services.Service = *service
	return service, nil
}

func (services TestServices) Delete(context context.Context, serviceNameOrId *string) error {
	*services.Deleted = *serviceNameOrId
	return nil
}

type FailServices struct {
}

func (streams FailServices) Create(context context.Context, service *gokong.Service) (*gokong.Service, error) {
	return service, errors.New("420 Enhance your calm")
}

func (services FailServices) Delete(context context.Context, serviceNameOrId *string) error {
	return nil
}

type ConflictServices struct {
}

func (streams ConflictServices) Create(context context.Context, service *gokong.Service) (*gokong.Service, error) {
	return service, errors.New("409 Conflict")
}

func (services ConflictServices) Delete(context context.Context, serviceNameOrId *string) error {
	return nil
}

type TestTargets struct {
	Created *[]gokong.Target
	Deleted *[]string
}

func (targets TestTargets) Create(context context.Context, target *gokong.Target) (*gokong.Target, error) {
	*targets.Created = append(*targets.Created, *target)
	return target, nil
}

func (targets TestTargets) Delete(context context.Context, upstreamNameOrId *string, targetOrId *string) error {
	*targets.Deleted = append(*targets.Deleted, *targetOrId)
	return nil
}

type FailTargets struct {
}

func (targets FailTargets) Create(context context.Context, target *gokong.Target) (*gokong.Target, error) {
	return target, errors.New("420 Enhance your calm")
}

func (targets FailTargets) Delete(context context.Context, upstreamNameOrId *string, targetOrId *string) error {
	return nil
}

type ConflictTargets struct {
}

func (targets ConflictTargets) Create(context context.Context, target *gokong.Target) (*gokong.Target, error) {
	return target, errors.New("409 Conflict")
}

func (targets ConflictTargets) Delete(context context.Context, upstreamNameOrId *string, targetOrId *string) error {
	return nil
}

type TestUpstreams struct {
	Created *gokong.Upstream
	Deleted *string
}

func (upstreams TestUpstreams) Create(context context.Context, upstream *gokong.Upstream) (*gokong.Upstream, error) {
	*upstreams.Created = *upstream
	return upstream, nil
}

func (upstreams TestUpstreams) Delete(context context.Context, upstreamNameOrId *string) error {
	*upstreams.Deleted = *upstreamNameOrId
	return nil
}

type FailUpstreams struct {
}

func (upstreams FailUpstreams) Create(context context.Context, upstream *gokong.Upstream) (*gokong.Upstream, error) {
	return upstream, errors.New("420 Enhance your calm")
}

func (upstreams FailUpstreams) Delete(context context.Context, upstreamNameOrId *string) error {
	return nil
}

type ConflictUpstreams struct {
}

func (upstreams ConflictUpstreams) Create(context context.Context, upstream *gokong.Upstream) (*gokong.Upstream, error) {
	return upstream, errors.New("409 Conflict")
}

func (upstreams ConflictUpstreams) Delete(context context.Context, upstreamNameOrId *string) error {
	return nil
}

/// ********************************************************************************************************************
/// TESTS

func TestRegistration_Register_CreateRouteCalled(t *testing.T) {
	mockClient, routes, service, _, _ := buildMockClient()

	registration, _ := NewRegistration(mockClient)
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

func TestRegistration_Register_CreateRouteIgnores409(t *testing.T) {
	mockClient, _, _, _, _ := buildMockClient()
	mockClient.Routes = ConflictRoutes{}

	registration, _ := NewRegistration(Client(mockClient))
	serviceDef := buildExampleServiceDef()

	err := registration.Register(serviceDef)
	if err != nil {
		t.Fatalf("Register should not have failed. 409 Conflict responses are ignored. %v", err)
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

func TestRegistration_Register_CreateServiceIgnores409(t *testing.T) {
	mockClient, _, _, _, _ := buildMockClient()
	mockClient.Services = ConflictServices{}

	registration, _ := NewRegistration(Client(mockClient))
	serviceDef := buildExampleServiceDef()

	err := registration.Register(serviceDef)
	if err != nil {
		t.Fatalf("Register should not have failed. 409 Conflict responses are ignored. %v", err)
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

func TestRegistration_Register_CreateTargetIgnores409(t *testing.T) {
	mockClient, _, _, _, _ := buildMockClient()
	mockClient.Targets = ConflictTargets{}

	registration, _ := NewRegistration(Client(mockClient))
	serviceDef := buildExampleServiceDef()

	err := registration.Register(serviceDef)
	if err != nil {
		t.Fatalf("Register should not have failed. 409 Conflict responses are ignored. %v", err)
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

func TestRegistration_Register_CreateUpstreamIgnores409(t *testing.T) {
	mockClient, _, _, _, _ := buildMockClient()
	mockClient.Upstreams = ConflictUpstreams{}

	registration, _ := NewRegistration(Client(mockClient))
	serviceDef := buildExampleServiceDef()

	err := registration.Register(serviceDef)
	if err != nil {
		t.Fatalf("Register should not have failed. 409 Conflict responses are ignored. %v", err)
	}
}

func TestRegistration_Deregister(t *testing.T) {
	client, routes, services, targets, upstreams := buildMockClient()
	registration, _ := NewRegistration(client)
	serviceDef := buildExampleServiceDef()

	err := registration.Deregister(serviceDef)

	if err != nil {
		t.Fatalf("Deregister failed with: %v", err)
	}

	if !reflect.DeepEqual(serviceDef.Targets, *targets.Deleted) {
		t.Fatal(fmt.Sprintf("Expected Targets.Delete to have been called with:\n\t%v, \nActual:\n\t%v", serviceDef.Targets, *targets.Deleted))
	}

	if serviceDef.UpstreamName != *upstreams.Deleted {
		t.Fatalf("Expected Upstreams.Delete to have been called with: %s, but was: %s", serviceDef.UpstreamName, *upstreams.Deleted)
	}

	expectedRoutes := []string{}
	for name, _ := range serviceDef.RoutesMap {
		expectedRoutes = append(expectedRoutes, name)
	}

	if !reflect.DeepEqual(expectedRoutes, *routes.Deleted) {
		t.Fatal(fmt.Sprintf("Expected Routes.Delete to have been called with:\n\t%v, \nActual:\n\t%v", expectedRoutes, *routes.Deleted))
	}

	if serviceDef.ServiceName != *services.Deleted {
		t.Fatalf(fmt.Sprintf("Expected Services.Delete to have been called with:\n\t%v, Actual: \n\t%v", serviceDef.ServiceName, *services.Deleted))
	}
}

/// ********************************************************************************************************************
/// HELPERS

func buildMockClient() (Client, TestRoutes, TestServices, TestTargets, TestUpstreams) {
	createdRoutes := []gokong.Route{}
	deletedRoutes := []string{}
	createdTargets := []gokong.Target{}
	deletedTargets := []string{}

	routes := TestRoutes{Created: &createdRoutes, Deleted: &deletedRoutes}
	services := TestServices{
		Service: new(gokong.Service),
		Deleted: new(string),
	}
	targets := TestTargets{Created: &createdTargets, Deleted: &deletedTargets}
	upstreams := TestUpstreams{Created: new(gokong.Upstream), Deleted: new(string)}

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
