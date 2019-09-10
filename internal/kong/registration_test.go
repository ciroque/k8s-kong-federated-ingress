package kong

import (
	"github.com/ciroque/k8s-kong-federated-ingress/internal/k8s"
	"testing"
)

type MockClient struct {

}

func (mockClient MockClient) CreateUpstream() error {
	return nil
}

func TestRegistration_Register(t *testing.T) {
	mockClient := new(MockClient)

	registration, _ := NewRegistration(mockClient)
	service := k8s.Service{
		Addresses: []string{},
		Name:      "test-service",
		Paths:     []string{},
		Port:      8080,
	}

	err := registration.Register(&service)
	if err != nil {
		t.Fatalf("Register failed with: %v", err)
	}

}
