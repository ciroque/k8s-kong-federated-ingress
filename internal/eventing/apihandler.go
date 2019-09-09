package eventing

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ciroque/k8s-kong-federated-ingress/internal/kong"
	networking "k8s.io/api/networking/v1beta1"
)

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

func (apiHandler *ApiHandler) Init() error {
	log.Info("ApiHandler.Init")
	return nil
}

func (apiHandler *ApiHandler) ObjectCreated(obj interface{}) error {
	log.Info("ApiHandler.ObjectCreated")
	ingress := obj.(*networking.Ingress)

	// THE NEW HOTNESS
	service, err := apiHandler.Translator.IngressToService(ingress)
	if err != nil {
		return fmt.Errorf("error handling ObjectCreated: %v", err)
	}
	_, _ = apiHandler.Registrar.Register(service)

	// OLD AND BUSTED
	return apiHandler.Kong.CreateKongObjects(ingress)
}

func (apiHandler *ApiHandler) ObjectDeleted(obj interface{}) error {
	log.Infof("ApiHandler.ObjectDeleted: %v", obj)
	ingress := obj.(*networking.Ingress)

	// THE NEW HOTNESS
	service, err := apiHandler.Translator.IngressToService(ingress)
	if err != nil {
		return fmt.Errorf("error handling ObjectDeleted: %v", err)
	}
	_ = apiHandler.Registrar.Deregister(service)

	// OLD AND BUSTED
	return apiHandler.Kong.DeleteKongObjects(ingress)
}

func (apiHandler *ApiHandler) ObjectUpdated(objOld, objNew interface{}) error {
	log.Info("ApiHandler.ObjectUpdated")
	oldIngress := objOld.(*networking.Ingress)
	newIngress := objNew.(*networking.Ingress)

	// THE NEW HOTNESS
	oldService, err := apiHandler.Translator.IngressToService(oldIngress)
	newService, err := apiHandler.Translator.IngressToService(newIngress)
	if err != nil {
		return fmt.Errorf("error handling ObjectUpdated: %v", err)
	}
	_, _ = apiHandler.Registrar.Modify(oldService, newService)

	// OLD AND BUSTED
	return apiHandler.Kong.UpdateKongObjects(oldIngress, newIngress)
}
