package client

import (
	"context"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetKubernetesClient returns Kubernetes client to use for the worker.
func GetKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func getConfig() (*rest.Config, error) {
	configPath := os.Getenv("KUBECONFIG")
	if configPath == "" {
		configPath = path.Join(os.Getenv("HOME"), ".kube/config")
	}
	if _, err := os.Stat(configPath); err == nil {
		logrus.WithField("path", configPath).Info("Using Kubernetes config based on config file")
		return clientcmd.BuildConfigFromFlags("", configPath)
	}
	logrus.Info("Using Kubernetes in-cluster config")
	return rest.InClusterConfig()
}

// Client performs requests to the Kubrentes API server.
type Client struct {
	clientSet *kubernetes.Clientset
}

// New returns a new instance of Client configured to connect to the API server.
func New() (*Client, error) {
	clientSet, err := GetKubernetesClient()
	if err != nil {
		return nil, err
	}
	return &Client{clientSet: clientSet}, nil
}

// GetMetrics requsts the /metrics endpoint from the server.
func (c *Client) GetMetrics(ctx context.Context) ([]byte, error) {
	return c.clientSet.Discovery().RESTClient().Get().AbsPath("/metrics").DoRaw(ctx)
}
