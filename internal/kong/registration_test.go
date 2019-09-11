package kong

import (
	"context"
	"github.com/ciroque/k8s-kong-federated-ingress/internal/k8s"
	gokong "github.com/hbagdi/go-kong/kong"
	"testing"
)

type MockClient struct {
	Upstreams UpstreamsInterface
}

type Upstreams struct {
	CreateCount *int // oddly this is the only way I could get this to work...
}

func (upstreams Upstreams) Create(context context.Context, upstream gokong.Upstream) (gokong.Upstream, error) {
	*upstreams.CreateCount = *(upstreams.CreateCount) + 1
	return upstream, nil
}

func TestRegistration_Register(t *testing.T) {
	upstreams := Upstreams{CreateCount: new(int)}
	mockClient := MockClient{
		Upstreams: upstreams,
	}

	registration, _ := NewRegistration(ClientInterface(mockClient))
	service := k8s.Service{
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
