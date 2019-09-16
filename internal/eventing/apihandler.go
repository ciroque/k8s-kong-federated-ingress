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

type K8s struct {
	Translator k8s.Translator
}

type Kong struct {
	Registrar  kong.Registrar
	Translator kong.Translator
}

type ApiHandler struct {
	K8s  K8s
	Kong Kong
}

func (apiHandler *ApiHandler) Init() error {
	log.Info("ApiHandler.Init")
	return nil
}

func (apiHandler *ApiHandler) ObjectCreated(obj interface{}) error {
	log.Info("ApiHandler.ObjectCreated")
	ingress := obj.(*networking.Ingress)

	serviceMap, err := apiHandler.K8s.Translator.IngressToService(ingress)
	if err != nil {
		return fmt.Errorf("error handling ObjectCreated: %v", err)
	}

	for serviceName, serviceDef := range serviceMap {
		kongService, err := apiHandler.Kong.Translator.ServiceToKong(serviceName, serviceDef)
		if err != nil {
			return fmt.Errorf("error translating service to Kong service: %v", err)
		}
		registrationErr := apiHandler.Kong.Registrar.Register(kongService)
		if registrationErr != nil {
			return fmt.Errorf("error registering Kong Service: %v", err)
		}
	}

	return nil
}

func (apiHandler *ApiHandler) ObjectDeleted(obj interface{}) error {
	log.Infof("ApiHandler.ObjectDeleted: %v", obj)
	ingress := obj.(*networking.Ingress)

	_, err := apiHandler.K8s.Translator.IngressToService(ingress)
	if err != nil {
		return fmt.Errorf("error handling ObjectDeleted: %v", err)
	}

	//return apiHandler.Registrar.Deregister(service)
	return nil
}

func (apiHandler *ApiHandler) ObjectUpdated(objOld, objNew interface{}) error {
	log.Info("ApiHandler.ObjectUpdated")
	oldIngress := objOld.(*networking.Ingress)
	newIngress := objNew.(*networking.Ingress)

	_, err := apiHandler.K8s.Translator.IngressToService(oldIngress)
	if err != nil {
		return fmt.Errorf("ObjectUpdated error translating oldIngress(%v): %v", oldIngress, err)
	}

	_, err2 := apiHandler.K8s.Translator.IngressToService(newIngress)
	if err2 != nil {
		return fmt.Errorf("ObjectUpdated error translating newIngress(%v): %v", newIngress, err)
	}

	//return apiHandler.Registrar.Modify(oldService, newService)
	return nil
}
