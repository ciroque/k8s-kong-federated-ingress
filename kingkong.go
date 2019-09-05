package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ciroque/kongo/client"
	v1 "k8s.io/api/core/v1"
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
	//log.Infof("KingKong Updating: %v", ingress)
	baseUrl := "http://localhost:8001"
	kongo, err := client.NewKongo(&baseUrl)
	if err != nil {
		return err
	}

	//addresses := []*string{}
	//for _, ingress := range ingress.Status.LoadBalancer.Ingress {
	//	address := fmt.Sprintf("%s:%i", ingress.IP, 80)
	//	addresses = append(addresses, &address)
	//}

	buildAddresses := func(backend networking.IngressBackend, ingresses []v1.LoadBalancerIngress) []*string {
		addresses := []*string{}
		for _, ingress := range ingresses {
			address := fmt.Sprintf("%s:%v", ingress.IP, backend.ServicePort.IntVal)
			addresses = append(addresses, &address)
		}
		return addresses
	}

	for ridx, rule := range ingress.Spec.Rules {
		for pidx, path := range rule.HTTP.Paths {
			addresses := buildAddresses(path.Backend, ingress.Status.LoadBalancer.Ingress)

			fmt.Println("---------------------------------------------------------------------------------------------")
			fmt.Println("Rule[", ridx, "]", "Path[", pidx, "] := ", fmt.Sprintf("{path: %s, serviceName: %s, servicePort: %v }", path.Path, path.Backend.ServiceName, path.Backend.ServicePort.IntVal))
			for _, address := range addresses {
				fmt.Println(fmt.Sprintf("Address: %v", *address))
			}
			fmt.Println("---------------------------------------------------------------------------------------------")

			k8sService := new(client.K8sService)

			k8sService.Name = path.Backend.ServiceName
			k8sService.Port = int(path.Backend.ServicePort.IntVal)
			k8sService.Path = path.Path /// NOTE: This could point to a subset of the path registered in the Ingress resource...
			k8sService.Addresses = addresses

			kongService, err := kongo.RegisterK8sService(k8sService)

			if err != nil {
				fmt.Println("ERROR Creating ", path.Backend.ServiceName, ": ", err)
			}

			fmt.Println(">>>> ", kongService)

		}
	}

	//address := "localhost"

	return err
}
