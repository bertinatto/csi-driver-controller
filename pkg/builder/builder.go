package builder

import (
	"os"

	apiextclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	dynamicclient "k8s.io/client-go/dynamic"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	osclientset "github.com/openshift/client-go/config/clientset/versioned"
)

// ClientBuilder can create a variety of kubernetes client interface with its embedded rest.Config.
type ClientBuilder struct {
	config *rest.Config
}

// APIExtClientOrDie returns the kubernetes client interface for extended kubernetes objects.
func (cb *ClientBuilder) APIExtClientOrDie(name string) apiextclient.Interface {
	return apiextclient.NewForConfigOrDie(rest.AddUserAgent(cb.config, name))
}

// KubeClientOrDie returns the kubernetes client interface for general kubernetes objects.
func (cb *ClientBuilder) KubeClientOrDie(name string) kubeclient.Interface {
	return kubeclient.NewForConfigOrDie(rest.AddUserAgent(cb.config, name))
}

// OSClientOrDie returns a client interface for general OpenShift objects.
func (cb *ClientBuilder) OSClientOrDie(name string) osclientset.Interface {
	return osclientset.NewForConfigOrDie(rest.AddUserAgent(cb.config, name))
}

// DynamicClientOrDie returns a dynamic client interface.
func (cb *ClientBuilder) DynamicClientOrDie(name string) dynamicclient.Interface {
	return dynamicclient.NewForConfigOrDie(rest.AddUserAgent(cb.config, name))
}

// NewClientBuilder returns a *ClientBuilder with the given kubeconfig.
func NewClientBuilder(kubeconfig string) (*ClientBuilder, error) {
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}

	if kubeconfig != "" {
		klog.Infof("Loading kube client config from path %q", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		klog.Infof("Using in-cluster kube client config")
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	return &ClientBuilder{
		config: config,
	}, nil
}
