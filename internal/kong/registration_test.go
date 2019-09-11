package kong

import (
	"context"
	gokong "github.com/hbagdi/go-kong/kong"
	"testing"
)

type MockClient struct {
	Services  ServicesInterface
	Upstreams UpstreamsInterface
}

type Services struct {
	CreateCount *int
}

func (streams Services) Create(ctx context.Context, service gokong.Service) (gokong.Service, error) {
	*streams.CreateCount = *streams.CreateCount + 1
	return service, nil
}

type Upstreams struct {
	CreateCount *int // oddly this is the only way I could get this to work...
}

func (upstreams Upstreams) Create(context context.Context, upstream gokong.Upstream) (gokong.Upstream, error) {
	*upstreams.CreateCount = *upstreams.CreateCount + 1
	return upstream, nil
}

func TestRegistration_Register_UpstreamCreated(t *testing.T) {
	mockClient, _, upstreams := buildMockClient()

	registration, _ := NewRegistration(ClientInterface(mockClient))
	service := ServiceDef{
		Addresses: []string{},
		Name:      "test-service",
		Paths:     []string{},
		Port:      8080,
	}

	err := registration.Register(service)
	if err != nil {
		t.Fatalf("Register failed with: %v", err)
	}

	if *upstreams.CreateCount != 1 {
		t.Fatal("Expected Upstreams.Create to have been called once. Actual count: ", *upstreams.CreateCount)
	}
}

func TestRegistration_Register_ServiceCreated(t *testing.T) {
	mockClient, services, _ := buildMockClient()

	registration, _ := NewRegistration(ClientInterface(mockClient))
	service := ServiceDef{
		Addresses: []string{},
		Name:      "test-service",
		Paths:     []string{},
		Port:      8080,
	}

	err := registration.Register(service)
	if err != nil {
		t.Fatalf("Register failed with: %v", err)
	}

	if *services.CreateCount != 1 {
		t.Fatal("Expected Services.Create to have been called once. Actual count: ", *services.CreateCount)
	}
}

func buildMockClient() (ClientInterface, Services, Upstreams) {
	services := Services{CreateCount: new(int)}
	upstreams := Upstreams{CreateCount: new(int)}
	mockClient := MockClient{
		Services:  services,
		Upstreams: upstreams,
	}
	return ClientInterface(mockClient), services, upstreams
}
