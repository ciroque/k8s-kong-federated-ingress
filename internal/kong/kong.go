/*

	DEPRECATING THIS ; THE NEW HOTNESS IS Translation / Registration implementation

*/

package kong

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ciroque/kongo/client"
	v1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
)

type Kong struct {
}

func (k *Kong) Init() error {
	return nil
}

func (k *Kong) CreateKongObjects(ingress *networking.Ingress) error {
	// TODO: Move this to the module level and use parameters to configure
	baseUrl := "http://localhost:8001"
	kongo, err := client.NewKongo(&baseUrl)
	if err != nil {
		return err
	}

	buildAddresses := func(backend networking.IngressBackend, ingresses []v1.LoadBalancerIngress) []*string {
		var addresses []*string
		for _, ingress := range ingresses {
			address := fmt.Sprintf("%s:%v", ingress.IP, backend.ServicePort.IntVal)
			addresses = append(addresses, &address)
		}
		return addresses
	}

	k8sService := new(client.K8sService)

	for _, rule := range ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			addresses := buildAddresses(path.Backend, ingress.Status.LoadBalancer.Ingress)

			k8sService.Name = path.Backend.ServiceName // TODO: Do we need to support multiple paths to the same service? If so this will need to change!
			k8sService.Port = int(path.Backend.ServicePort.IntVal)
			k8sService.Path = path.Path /// NOTE: This could point to a subset of the path registered in the Ingress resource...
			k8sService.Addresses = addresses

			_, err := kongo.RegisterK8sService(k8sService)
			if err != nil {
				fmt.Println("ERROR Creating ", path.Backend.ServiceName, ": ", err)
			}
		}
	}
	return nil
}

func (k *Kong) DeleteKongObjects(ingress *networking.Ingress) error {
	log.Infof("Kong Deleting: %v", ingress)

	var gerr error

	// TODO: Move this to the module level and use parameters to configure
	baseUrl := "http://localhost:8001"
	kongo, err := client.NewKongo(&baseUrl)
	if err != nil {
		return err
	}

	k8sService := new(client.K8sService)

	for _, rule := range ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			k8sService.Name = path.Backend.ServiceName
			err := kongo.DeregisterK8sService(k8sService)
			if err != nil {
				gerr = fmt.Errorf("error deregistering service '%s': %v", k8sService.Name, err)
			}
		}
	}

	return gerr
}

func (k *Kong) UpdateKongObjects(oldIngress *networking.Ingress, newIngress *networking.Ingress) error {
	log.Infof("Kong Updating: %v => %v", oldIngress, newIngress)
	//baseUrl := "http://localhost:8001"
	//kongo, err := client.NewKongo(&baseUrl)
	//if err != nil {
	//	return err
	//}
	return nil
}
