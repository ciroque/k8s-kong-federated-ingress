package eventing

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.marchex.com/marchex/k8s-kong-federated-ingress/internal/k8s"
	"github.marchex.com/marchex/k8s-kong-federated-ingress/internal/kong"
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
		return fmt.Errorf("ApiHandler::ObjectCreated error handling ObjectCreated: %#v", err)
	}

	for serviceName, serviceDef := range serviceMap {
		kongService, err := apiHandler.Kong.Translator.ServiceToKong(serviceName, serviceDef)
		if err != nil {
			return fmt.Errorf("ApiHandler::ObjectCreated error translating service to Kong service: %#v", err)
		}
		err = apiHandler.Kong.Registrar.Register(kongService)
		if err != nil {
			return fmt.Errorf("ApiHandler::ObjectCreated error registering Kong Service: %#v", err)
		}
	}

	return nil
}

func (apiHandler *ApiHandler) ObjectDeleted(obj interface{}) error {
	log.Infof("ApiHandler.ObjectDeleted: %#v", obj)
	ingress := obj.(*networking.Ingress)

	serviceMap, err := apiHandler.K8s.Translator.IngressToService(ingress)
	if err != nil {
		return fmt.Errorf("error handling ObjectDeleted: %#v", err)
	}

	for serviceName, serviceDef := range serviceMap {
		kongService, err := apiHandler.Kong.Translator.ServiceToKong(serviceName, serviceDef)
		if err != nil {
			return fmt.Errorf("ApiHandler::ObjectDeleted failed to translate a k8s.ServiceDef to a kong.ServiceDef: %#v", err)
		}

		err = apiHandler.Kong.Registrar.Deregister(kongService)
		if err != nil {
			return fmt.Errorf("ApiHandler::ObjectDeleted failed to Deregister a service: %#v. Error: %#v", kongService, err)
		}
	}

	return nil
}

func (apiHandler *ApiHandler) ObjectUpdated(originalResource, revisedResource interface{}) error {
	log.Info("ApiHandler.ObjectUpdated")
	previousIngress := originalResource.(*networking.Ingress)
	revisedIngress := revisedResource.(*networking.Ingress)

	originalServiceMap, err := apiHandler.K8s.Translator.IngressToService(previousIngress)
	if err != nil {
		return fmt.Errorf("ObjectUpdated error translating previousIngress(%#v): %#v", previousIngress, err)
	}

	revisedServiceMap, err := apiHandler.K8s.Translator.IngressToService(revisedIngress)
	if err != nil {
		return fmt.Errorf("ObjectUpdated error translating revisedIngress(%#v): %#v", revisedIngress, err)
	}

	var gerr error

	for revisedServiceName, revisedServiceDef := range revisedServiceMap {
		if _, found := originalServiceMap[revisedServiceName]; found {
			if kongService, err := apiHandler.Kong.Translator.ServiceToKong(revisedServiceName, revisedServiceDef); err == nil {
				if apiHandler.Kong.Registrar.Register(kongService) != nil {
					gerr = fmt.Errorf("ApiHandler::ObjectUpdated failed to Register a service: %#v. Error: %#v. Previous errors: %v", kongService, err, gerr)
				}
			}
		}
	}

	return gerr
}
