package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ciroque/kongo/kongClient"
	networking "k8s.io/api/networking/v1beta1"
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

func (k *KingKong) UpsertKongObjects(ingress *networking.Ingress) error {
	log.Infof("KingKong Updating: %v", ingress)
	baseUrl := "http://localhost:8001"
	kongo, err := kongClient.NewKongo(&baseUrl)
	if err != nil {
		return err
	}

	//address := "localhost"

	k8sService := new(kongClient.K8sService)

	k8sService.Name = ingress.Name
	//k8sService.Port = int(ingress.Backend.ServicePort.IntVal)
	//k8sService.Path = ingress.Status.  /// Could be pinned to the current static DNS entry for Kong (King)
	//k8sService.Addresses = []*string{&address}

	kongService, err := kongo.RegisterK8sService(k8sService)

	fmt.Println(">>>> ", kongService)

	return err
}
