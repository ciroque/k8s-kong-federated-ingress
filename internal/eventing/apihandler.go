package eventing

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ciroque/k8s-kong-federated-ingress/internal/kong"
	networking "k8s.io/api/networking/v1beta1"
)

// Handler interface contains the methods that are required
type Handler interface {
	Init() error
	ObjectCreated(obj interface{}) error
	ObjectDeleted(obj interface{}) error
	ObjectUpdated(objOld, objNew interface{}) error
}

type ApiHandler struct {
	Kong       kong.Kong
	Translator kong.Translator /// STEVE: This can be mocked for tests
	Registrar  kong.Registrar  /// STEVE: This can be mocked for tests
}

func (t *ApiHandler) Init() error {
	log.Info("ApiHandler.Init")
	return nil
}

func (t *ApiHandler) ObjectCreated(obj interface{}) error {
	log.Info("ApiHandler.ObjectCreated")
	ingress := obj.(*networking.Ingress)

	// THE NEW HOTNESS
	service, err := t.Translator.IngressToService(ingress)
	if err != nil {
		return fmt.Errorf("error handling ObjectCreated: %v", err)
	}
	_, _ = t.Registrar.Register(service)

	// OLD AND BUSTED
	return t.Kong.CreateKongObjects(ingress)
}

func (t *ApiHandler) ObjectDeleted(obj interface{}) error {
	log.Infof("ApiHandler.ObjectDeleted: %v", obj)
	ingress := obj.(*networking.Ingress)

	// THE NEW HOTNESS
	service, err := t.Translator.IngressToService(ingress)
	if err != nil {
		return fmt.Errorf("error handling ObjectDeleted: %v", err)
	}
	_ = t.Registrar.Deregister(service)

	// OLD AND BUSTED
	return t.Kong.DeleteKongObjects(ingress)
}

func (t *ApiHandler) ObjectUpdated(objOld, objNew interface{}) error {
	log.Info("ApiHandler.ObjectUpdated")
	oldIngress := objOld.(*networking.Ingress)
	newIngress := objNew.(*networking.Ingress)

	// THE NEW HOTNESS
	oldService, err := t.Translator.IngressToService(oldIngress)
	newService, err := t.Translator.IngressToService(newIngress)
	if err != nil {
		return fmt.Errorf("error handling ObjectUpdated: %v", err)
	}
	_, _ = t.Registrar.Modify(oldService, newService)

	// OLD AND BUSTED
	return t.Kong.UpdateKongObjects(oldIngress, newIngress)
}
