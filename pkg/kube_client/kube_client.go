package kube_client

import (
	// metavif1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	loadingRules *clientcmd.ClientConfigLoadingRules
	kubeConfig   clientcmd.ClientConfig
	config       *restclient.Config
	restConfig   *restclient.Config
	clientset    *kubernetes.Clientset
)

func LoadingRules() *clientcmd.ClientConfigLoadingRules {
	if loadingRules == nil {
		loadingRules = clientcmd.NewDefaultClientConfigLoadingRules()
	}

	return loadingRules
}

func KubeConfig() clientcmd.ClientConfig {
	if kubeConfig == nil {
		configOverrides := &clientcmd.ConfigOverrides{}
		kubeConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(LoadingRules(), configOverrides)
	}

	return kubeConfig
}

func RestConfig() (c *restclient.Config, err error) {
	if restConfig == nil {
		c, err = kubeConfig.ClientConfig()
	}
	return c, err
}

func ClientSet() (cs *kubernetes.Clientset, err error) {
	if restConfig, err = RestConfig(); err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(restConfig)
}
