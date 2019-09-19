/*

	Entrypoint for k8s kong federated ingress controller.

	Much of this code was lifted from: https://github.com/trstringer/k8s-controller-core-resource

*/

package main

import (
	"crypto/tls"
	"errors"
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

func GetKubernetesClient(config *Config) kubernetes.Interface {
	k8sConfig, err := clientcmd.BuildConfigFromFlags("", config.KubeConfigPath)
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
	KubeConfigPath string
	KongHost       string
}

/// Pulling from environment variables for now. TODO: Use Consul to home the configuration values (optionally)
func NewConfig() (*Config, error) {
	config := new(Config)

	config.KongHost = os.Getenv("KONG_HOST")
	if config.KongHost == "" {
		return nil, errors.New("the KONG_HOST variable is not defined. This is required")
	}

	config.KubeConfigPath = os.Getenv("KUBE_CONFIG_FILE")
	if config.KubeConfigPath == "" {
		return nil, errors.New("the KUBE_CONFIG_FILE variable is not defined. This is required")
	}

	return config, nil
}

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("Failed to load Config: %#v", err)
	}

	client := GetKubernetesClient(config)

	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return client.NetworkingV1beta1().Ingresses(meta_v1.NamespaceAll).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return client.NetworkingV1beta1().Ingresses(meta_v1.NamespaceAll).Watch(options)
			},
		},
		&networking.Ingress{},
		0, // no resync (period of 0)
		cache.Indexers{},
	)

	eventQueue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

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
		log.Fatalf("Unable to create a Kong Client: %#v\n", err)
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

	stopCh := make(chan struct{})
	defer close(stopCh)

	go controller.Run(stopCh)

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)
	<-sigTerm
}

func buildKongClient(client *http.Client, config *Config) (*gokong.Client, error) {
	return gokong.NewClient(&config.KongHost, client)
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
