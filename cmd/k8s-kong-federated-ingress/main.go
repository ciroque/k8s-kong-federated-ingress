/*

	Entrypoint for k8s kong federated ingress controller.

	Much of this code was lifted from: https://github.com/trstringer/k8s-controller-core-resource

*/

package main

import (
	"github.com/ciroque/k8s-kong-federated-ingress/internal/eventing"
	"github.com/ciroque/k8s-kong-federated-ingress/internal/kong"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"

	networking "k8s.io/api/networking/v1beta1"
	//api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

// retrieve the Kubernetes cluster client from outside of the cluster
func GetKubernetesClient() kubernetes.Interface {
	// construct the path to resolve to `~/.kube/config`
	kubeConfigPath := os.Getenv("HOME") + "/.kube/config"

	// create the config from the path
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	// generate the client based off of the config
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	log.Info("Successfully constructed k8s client")
	return client
}

func main() {

	// get the Kubernetes client for connectivity
	client := GetKubernetesClient()

	// create the informer so that we can not only list resources
	// but also watch them for all Ingress resources in the default namespace
	informer := cache.NewSharedIndexInformer(
		// the ListWatch contains two different functions that our
		// informer requires: ListFunc to take care of listing and watching
		// the resources we want to handle
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				// list all of the ingresses (Ingress resource) in the default namespace
				//return client.NetworkingV1beta1().Ingresses("app-services").List(options)
				return client.NetworkingV1beta1().Ingresses(meta_v1.NamespaceAll).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				//return client.NetworkingV1beta1().Ingresses("app-services").Watch(options)
				return client.NetworkingV1beta1().Ingresses(meta_v1.NamespaceAll).Watch(options)
			},
		},
		&networking.Ingress{},
		0, // no resync (period of 0)
		cache.Indexers{},
	)

	eventQueue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// add event handlers to handle the three types of events for resources:
	//  - adding new resources
	//  - updating existing resources
	//  - deleting resources
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			e := eventing.NewEvent(eventing.Created, obj, nil)
			eventQueue.Add(e)
			log.Infof("Added Created event to eventQueue %v", e)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			e := eventing.NewEvent(eventing.Updated, newObj, oldObj)
			eventQueue.Add(e)
			log.Infof("Added Updated event to eventQueue %v", e)
		},
		DeleteFunc: func(obj interface{}) {
			e := eventing.NewEvent(eventing.Deleted, obj, nil)
			eventQueue.Add(e)
			log.Infof("Added Deleted event to eventQueue %v", e)
		},
	})

	// construct the Controller object which has all of the necessary components to
	// handle logging, connections, informing (listing and watching), the queue,
	// and the handler
	controller := eventing.Controller{
		Logger:    log.NewEntry(log.New()),
		Clientset: client,
		Informer:  informer,
		Queue:     eventQueue,
		Handler: eventing.ApiHandler{
			Translator: &kong.Translation{},
			Registrar:  &kong.Registration{},
		},
	}

	// use a channel to synchronize the finalization for a graceful shutdown
	stopCh := make(chan struct{})
	defer close(stopCh)

	// run the controller loop to process items
	go controller.Run(stopCh)

	// use a channel to handle OS signals to terminate and gracefully shut
	// down processing
	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)
	<-sigTerm
}
