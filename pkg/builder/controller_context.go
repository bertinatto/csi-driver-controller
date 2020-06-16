package builder

import (
	"math/rand"
	"time"

	apiextinformers "k8s.io/apiextensions-apiserver/pkg/client/informers/externalversions"
	"k8s.io/client-go/informers"

	osinformers "github.com/openshift/client-go/config/informers/externalversions"
)

const (
	minResyncPeriod = 20 * time.Minute
)

func resyncPeriod() func() time.Duration {
	return func() time.Duration {
		factor := rand.Float64() + 1
		return time.Duration(float64(minResyncPeriod.Nanoseconds()) * factor)
	}
}

// ControllerContext stores all the informers for a variety of kubernetes objects.
type ControllerContext struct {
	ClientBuilder                 *ClientBuilder
	APIExtInformerFactory         apiextinformers.SharedInformerFactory
	KubeNamespacedInformerFactory informers.SharedInformerFactory
	OSInformerFactory             osinformers.SharedInformerFactory

	Stop <-chan struct{}

	InformersStarted chan struct{}

	ResyncPeriod func() time.Duration
}

// NewControllerContext creates the ControllerContext with the ClientBuilder.
func NewControllerContext(cb *ClientBuilder, stop <-chan struct{}, operandNamespace string) *ControllerContext {
	apiExtClient := cb.APIExtClientOrDie("apiext-shared-informer")
	apiExtSharedInformer := apiextinformers.NewSharedInformerFactoryWithOptions(apiExtClient, resyncPeriod()(),
		apiextinformers.WithNamespace(operandNamespace))

	kubeClient := cb.KubeClientOrDie("kube-shared-informer")
	kubeNamespacedSharedInformer := informers.NewFilteredSharedInformerFactory(kubeClient, resyncPeriod()(), operandNamespace, nil)

	osClient := cb.OSClientOrDie("openshift-shared-informer")
	osInformerFactory := osinformers.NewSharedInformerFactory(osClient, resyncPeriod()())

	return &ControllerContext{
		ClientBuilder:                 cb,
		APIExtInformerFactory:         apiExtSharedInformer,
		KubeNamespacedInformerFactory: kubeNamespacedSharedInformer,
		OSInformerFactory:             osInformerFactory,
		Stop:                          stop,
		InformersStarted:              make(chan struct{}),
		ResyncPeriod:                  resyncPeriod(),
	}
}
