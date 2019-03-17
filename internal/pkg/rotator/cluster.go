package rotator

import (
	"context"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", "/tmp/.kube/config") // TODO
	if err != nil {
		return nil, err
	}
	kc, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kc, nil
}

func DescribeClusterNodes(kc kubernetes.Clientset) (*sets.String, error) {
	return nil, nil
}

func AwaitNewNodeReady(ctx context.Context, kc kubernetes.Clientset, nodes *sets.String) error {
	return nil
}

func DrainNodeByProviderID(ctx context.Context, kc kubernetes.Clientset, providerID string) error {
	return nil
}

func CordonNodeByProviderID(ctx context.Context, kc kubernetes.Clientset, providerID string) error {
	return nil
}
