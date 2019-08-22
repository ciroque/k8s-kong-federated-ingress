package main

import (
	log "github.com/Sirupsen/logrus"
	networkv1beta1 "k8s.io/api/networking/v1beta1"
)

type KingKong struct {
}

// Init handles any handler initialization
func (k *KingKong) Init() error {
	log.Info("TestHandler.Init")
	return nil
}

func (k *KingKong) DeleteKongObjects(ingressName string) {
	log.Infof("KingKong Deleting: %s", ingressName)
}

func (k *KingKong) UpsertKongObjects(ingressSpec networkv1beta1.IngressSpec, kongHosts string) {
	log.Infof("KingKong Updating: %s \n\n %s", ingressSpec, kongHosts)
}
