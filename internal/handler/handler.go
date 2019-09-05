package handler

import (
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

// K8sApiHandler is a sample implementation of Handler
type K8sApiHandler struct {
	Kong kong.Kong
}

// Init handles any Handler initialization
func (t *K8sApiHandler) Init() error {
	log.Info("K8sApiHandler.Init")
	return nil
}

// ObjectCreated is called when an object is created
func (t *K8sApiHandler) ObjectCreated(obj interface{}) error {
	log.Info("K8sApiHandler.ObjectCreated")
	ingress := obj.(*networking.Ingress)
	return t.Kong.CreateKongObjects(ingress)
}

// ObjectDeleted is called when an object is deleted
func (t *K8sApiHandler) ObjectDeleted(obj interface{}) error {
	log.Infof("K8sApiHandler.ObjectDeleted: %v", obj)
	ingress := obj.(*networking.Ingress)
	return t.Kong.DeleteKongObjects(ingress)
}

// ObjectUpdated is called when an object is updated
func (t *K8sApiHandler) ObjectUpdated(objOld, objNew interface{}) error {
	log.Info("K8sApiHandler.ObjectUpdated")
	oldIngress := objOld.(*networking.Ingress)
	newIngress := objNew.(*networking.Ingress)
	return t.Kong.UpdateKongObjects(oldIngress, newIngress)
}
