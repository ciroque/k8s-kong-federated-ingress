package eventing

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	Logger    *log.Entry
	Clientset kubernetes.Interface
	Queue     workqueue.RateLimitingInterface
	Informer  cache.SharedIndexInformer
	Handler   ApiHandler
}

func (c *Controller) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.Queue.ShutDown()
	c.Logger.Info("Controller.Run: initiating")
	go c.Informer.Run(stopCh)
	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Error syncing cache"))
		return
	}
	c.Logger.Info("Controller.Run: cache sync complete")
	wait.Until(c.runWorker, time.Second, stopCh)
}

// HasSynced allows us to satisfy the Controller interface
// by wiring up the Informer's HasSynced method to it
func (c *Controller) HasSynced() bool {
	return c.Informer.HasSynced()
}

// runWorker executes the loop to process new items added to the Queue
func (c *Controller) runWorker() {
	log.Info("Controller.runWorker: starting")

	// invoke processNextItem to fetch and consume the next change
	// to a watched or listed resource
	for c.processNextItem() {
		log.Info("Controller.runWorker: processing next item")
	}

	log.Info("Controller.runWorker: completed")
}

func (c *Controller) processNextItem() bool {
	log.Info("Controller.processNextItem: start")

	e, quit := c.Queue.Get()

	if quit {
		return false
	}

	defer c.Queue.Done(e)

	withRetry := func(err error) {
		if err != nil {
			if c.Queue.NumRequeues(e) < 5 {
				log.Warnf("Requeuing %#v due to error: %#v", e, err)
				c.Queue.AddRateLimited(e)
			} else {
				log.Errorf("Dropping %#v due to too many errors. Latest error: %#v", e, err)
				c.Queue.Forget(e)
			}
		}
	}

	event := e.(Event)

	switch event.Type {
	case Created:
		{
			withRetry(c.Handler.ObjectCreated(event.Resource))
		}
	case Deleted:
		{
			withRetry(c.Handler.ObjectDeleted(event.Resource))
		}
	case Updated:
		{
			withRetry(c.Handler.ObjectUpdated(event.PreviousResource, event.Resource))
		}
	}

	// keep the worker loop running by returning true
	return true
}
