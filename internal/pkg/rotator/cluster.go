package rotator

import (
	"context"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	DefaultNodeAwaitJoinTimeout      = 30 * time.Second
	DefaultNodeAwaitReadinessTimeout = 10 * time.Second
)

func NewKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", "/tmp/.kube/config") // TODO
	if err != nil {
		return nil, err
	}
	k8s, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return k8s, nil
}

func DescribeClusterNodes(k8s *kubernetes.Clientset) (sets.String, error) {
	list, err := k8s.CoreV1().Nodes().List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	nodes := sets.NewString()
	for _, node := range list.Items {
		nodes.Insert(string(node.UID))
	}
	return nodes, nil
}

func AwaitNewNodeReady(ctx context.Context, k8s *kubernetes.Clientset, nodes sets.String) error {
	errors := make(chan error)
	go func() { errors <- awaitNewNodeReady(k8s, nodes) }()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errors:
		return err
	}
}

func awaitNewNodeReady(k8s *kubernetes.Clientset, nodes sets.String) error {
	node, err := awaitNewNodeJoin(k8s, nodes)
	if err != nil {
		return err
	}
	if err := awaitNodeReadiness(k8s, node); err != nil {
		return err
	}
	return nil
}

func awaitNewNodeJoin(k8s *kubernetes.Clientset, known sets.String) (*corev1.Node, error) {
	for {
		log.Printf("Waiting %s for new node to join cluster...", DefaultNodeAwaitJoinTimeout.String())
		time.Sleep(DefaultNodeAwaitJoinTimeout)

		list, err := k8s.CoreV1().Nodes().List(v1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for _, node := range list.Items {
			if known.Has(string(node.UID)) {
				continue
			}
			log.Printf("Node '%s' joined cluster.", node.Name)
			return &node, nil
		}
	}
}

func awaitNodeReadiness(k8s *kubernetes.Clientset, node *corev1.Node) error {
	for {
		log.Printf("Waiting %s for new node to be ready...", DefaultNodeAwaitReadinessTimeout.String())
		time.Sleep(DefaultNodeAwaitReadinessTimeout)

		n, err := k8s.CoreV1().Nodes().Get(node.Name, v1.GetOptions{})
		if err != nil {
			return err
		}
		for _, c := range n.Status.Conditions {
			if c.Type == corev1.NodeReady && c.Status == corev1.ConditionTrue {
				return nil
			}
		}
	}
}

func DrainNodeByInstanceID(ctx context.Context, k8s *kubernetes.Clientset, id string) error {
	errors := make(chan error)
	go func() { errors <- drainNodeByInstanceID(k8s, id) }()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errors:
		return err
	}
}

func drainNodeByInstanceID(k8s *kubernetes.Clientset, id string) error {
	return nil
}

func CordonNodeByProviderID(ctx context.Context, k8s *kubernetes.Clientset, providerID string) error {
	return nil
}
