package eventing

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ciroque/k8s-kong-federated-ingress/internal/k8s"
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
	Translator k8s.Translator
	Registrar  kong.Registrar
}

func (apiHandler *ApiHandler) Init() error {
	log.Info("ApiHandler.Init")
	return nil
}

func (apiHandler *ApiHandler) ObjectCreated(obj interface{}) error {
	log.Info("ApiHandler.ObjectCreated")
	ingress := obj.(*networking.Ingress)

	service, err := apiHandler.Translator.IngressToK8sService(ingress)
	if err != nil {
		return fmt.Errorf("error handling ObjectCreated: %v", err)
	}

	return apiHandler.Registrar.Register(service)
}

func (apiHandler *ApiHandler) ObjectDeleted(obj interface{}) error {
	log.Infof("ApiHandler.ObjectDeleted: %v", obj)
	ingress := obj.(*networking.Ingress)

	service, err := apiHandler.Translator.IngressToK8sService(ingress)
	if err != nil {
		return fmt.Errorf("error handling ObjectDeleted: %v", err)
	}

	return apiHandler.Registrar.Deregister(service)
}

func (apiHandler *ApiHandler) ObjectUpdated(objOld, objNew interface{}) error {
	log.Info("ApiHandler.ObjectUpdated")
	oldIngress := objOld.(*networking.Ingress)
	newIngress := objNew.(*networking.Ingress)

	oldService, err := apiHandler.Translator.IngressToK8sService(oldIngress)
	if err != nil {
		return fmt.Errorf("ObjectUpdated error translating oldIngress(%v): %v", oldIngress, err)
	}

	newService, err := apiHandler.Translator.IngressToK8sService(newIngress)
	if err != nil {
		return fmt.Errorf("ObjectUpdated error translating newIngress(%v): %v", newIngress, err)
	}

	return apiHandler.Registrar.Modify(oldService, newService)
}
