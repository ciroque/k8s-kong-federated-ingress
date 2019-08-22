package main

import (
	log "github.com/Sirupsen/logrus"
	//core_v1 "k8s.io/api/core/v1"
	network_v1beta1 "k8s.io/api/networking/v1beta1"
)

// Handler interface contains the methods that are required
type Handler interface {
	Init() error
	ObjectCreated(obj interface{})
	ObjectDeleted(key interface{}, obj interface{})
	ObjectUpdated(objOld, objNew interface{})
}

// TestHandler is a sample implementation of Handler
type TestHandler struct{
	kingKong KingKong

}

// Init handles any handler initialization
func (t *TestHandler) Init() error {
	log.Info("TestHandler.Init")
	return nil
}

// ObjectCreated is called when an object is created
func (t *TestHandler) ObjectCreated(obj interface{}) {
	log.Info("TestHandler.ObjectCreated")
	// assert the type to a Ingress object to pull out relevant data
	ingress := obj.(*network_v1beta1.Ingress)

	ingressSpec := ingress.Spec
	kongHosts := ingress.Annotations["ingress.marchex.net/kong-hosts"]
	t.kingKong.UpsertKongObjects(ingressSpec, kongHosts)
	//
	//log.Infof("Name: %s", ingress.Name)
	//log.Infof("String: %s", ingress.Spec)
	//log.Infof("Kong Hosts to update: %s", ingress.Annotations["ingress.marchex.net/kong-hosts"])
}

// ObjectDeleted is called when an object is deleted
func (t *TestHandler) ObjectDeleted(key interface{}, obj interface{}) {
	log.Infof("TestHandler.ObjectDeleted: %s", key)
	t.kingKong.DeleteKongObjects(key.(string))
}

// ObjectUpdated is called when an object is updated
func (t *TestHandler) ObjectUpdated(objOld, objNew interface{}) {
	log.Info("TestHandler.ObjectUpdated")
}
