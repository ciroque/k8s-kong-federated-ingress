/*

	Entrypoint for k8s kong federated ingress controller.

	Much of this code was lifted from: https://github.com/trstringer/k8s-controller-core-resource

*/

package main

import (
	"crypto/tls"
	log "github.com/Sirupsen/logrus"
	"github.com/ciroque/k8s-kong-federated-ingress/internal/eventing"
	"github.com/ciroque/k8s-kong-federated-ingress/internal/k8s"
	"github.com/ciroque/k8s-kong-federated-ingress/internal/kong"
	gokong "github.com/hbagdi/go-kong/kong"
	networking "k8s.io/api/networking/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

func GetKubernetesClient(config Config) kubernetes.Interface {
	kubeConfigPath := os.Getenv(config.KubeConfigPathEnvironmentVariable)

	k8sConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Fatalf("GetKubernetesClient: %v", err)
	}

	client, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		log.Fatalf("GetKubernetesClient: %#v", err)
	}

	log.Info("Successfully constructed k8s client")
	return client
}

type Config struct {
	KongHostEnvironmentVariable       string
	KubeConfigPathEnvironmentVariable string
}

/// Pulling from environment variables for now. TODO: Use Consul to home the configuration values (optionally)
func NewConfig() Config {
	return Config{
		KongHostEnvironmentVariable:       "KONG_HOST",
		KubeConfigPathEnvironmentVariable: "KUBE_CONFIG_FILE",
	}
}

func main() {
	config := NewConfig()

	// get the Kubernetes client for connectivity
	client := GetKubernetesClient(config)

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
			log.Infof("Added Created event to eventQueue %#v", e)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			e := eventing.NewEvent(eventing.Updated, newObj, oldObj)
			eventQueue.Add(e)
			log.Infof("Added Updated event to eventQueue %#v", e)
		},
		DeleteFunc: func(obj interface{}) {
			e := eventing.NewEvent(eventing.Deleted, obj, nil)
			eventQueue.Add(e)
			log.Infof("Added Deleted event to eventQueue %#v", e)
		},
	})

	httpClient := buildHttpClient()
	kongClient, err := buildKongClient(httpClient, config)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "Unable to create a Kong Client (Ensure that the %s environment variable is set!): %#v\n", config.KongHostEnvironmentVariable, err)
		//os.Exit(1)
		log.Fatalf("Unable to create a Kong Client (Ensure that the %s environment variable is set!): %#v\n", config.KongHostEnvironmentVariable, err)
	}

	eventingK8s := eventing.K8s{Translator: &k8s.Translation{}}
	eventingKong := eventing.Kong{
		Registrar: &kong.Registration{
			Kong: kong.Client{
				Routes:    kong.Routes{Kong: *kongClient},
				Services:  kong.Services{Kong: *kongClient},
				Targets:   kong.Targets{Kong: *kongClient},
				Upstreams: kong.Upstreams{Kong: *kongClient},
			},
		},
		Translator: &kong.Translation{},
	}

	controller := eventing.Controller{
		Logger:    log.NewEntry(log.New()),
		Clientset: client,
		Informer:  informer,
		Queue:     eventQueue,
		Handler: eventing.ApiHandler{
			K8s:  eventingK8s,
			Kong: eventingKong,
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

func buildKongClient(client *http.Client, config Config) (*gokong.Client, error) {
	host := os.Getenv(config.KongHostEnvironmentVariable)
	return gokong.NewClient(&host, client)
}

func buildHttpClient() *http.Client {
	headers := []string{"Content-Type: application/json", "Accept: application/json"}
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = true

	transport := http.DefaultTransport.(*http.Transport)

	transport.TLSClientConfig = &tlsConfig

	httpClient := http.Client{
		Transport:     transport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 10,
	}

	httpClient.Transport = &eventing.RoundTripper{
		Headers:      headers,
		RoundTripper: transport,
	}

	return &httpClient
}
