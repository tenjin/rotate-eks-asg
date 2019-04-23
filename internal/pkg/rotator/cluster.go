package rotator

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/tenjin/rotate-eks-asg/internal/pkg/cmd"
)

const (
	DefaultKubeConfigPath = "/tmp/.kube/config"
)

var (
	DefaultNodeAwaitJoinTimeout      = 30 * time.Second
	DefaultNodeAwaitReadinessTimeout = 10 * time.Second
)

func NewKubernetesClient() (*kubernetes.Clientset, error) {
	kcfg := os.Getenv("KUBECONFIG")
	if kcfg == "" {
		kcfg = DefaultKubeConfigPath
	}

	config, err := clientcmd.BuildConfigFromFlags("", kcfg)
	if err != nil {
		return nil, err
	}
	k8s, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return k8s, nil
}

func GetClusterNodeSet(k8s *kubernetes.Clientset) (sets.String, error) {
	nodes, err := getClusterNodes(k8s)
	if err != nil {
		return nil, err
	}
	set := sets.NewString()
	for _, node := range nodes {
		set.Insert(string(node.UID))
	}
	return set, nil
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

func awaitNewNodeJoin(k8s *kubernetes.Clientset, known sets.String) (*coreV1.Node, error) {
	for {
		log.Printf("Waiting %s for new node to join cluster...", DefaultNodeAwaitJoinTimeout.String())
		time.Sleep(DefaultNodeAwaitJoinTimeout)

		nodes, err := getClusterNodes(k8s)
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			if known.Has(string(node.UID)) {
				continue
			}
			log.Printf("Node '%s' joined cluster.", node.Name)
			return node, nil
		}
	}
}

func awaitNodeReadiness(k8s *kubernetes.Clientset, node *coreV1.Node) error {
	for {
		log.Printf("Waiting %s for new node to be ready...", DefaultNodeAwaitReadinessTimeout.String())
		time.Sleep(DefaultNodeAwaitReadinessTimeout)

		n, err := k8s.CoreV1().Nodes().Get(node.Name, v1.GetOptions{})
		if err != nil {
			return err
		}
		for _, c := range n.Status.Conditions {
			if c.Type == coreV1.NodeReady && c.Status == coreV1.ConditionTrue {
				return nil
			}
		}
	}
}

func GetNodeNameByInstanceID(k8s *kubernetes.Clientset, id string) (string, error) {
	nodes, err := getClusterNodes(k8s)
	if err != nil {
		return "", err
	}
	for _, node := range nodes {
		if strings.HasSuffix(node.Spec.ProviderID, id) {
			return node.Name, nil
		}
	}
	return "", fmt.Errorf("node '%s' is not part of the cluster", id)
}

func getClusterNodes(k8s *kubernetes.Clientset) ([]*coreV1.Node, error) {
	list, err := k8s.CoreV1().Nodes().List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	nodes := make([]*coreV1.Node, 0, len(list.Items))
	for _, node := range list.Items {
		n := node
		nodes = append(nodes, &n)
	}
	return nodes, nil
}

func DrainNodeByName(ctx context.Context, name string) error {
	return kubectl(ctx, "drain", "--delete-local-data=true", " --ignore-daemonsets=true", name)
}

func CordonNodeByName(ctx context.Context, name string) error {
	return kubectl(ctx, "cordon", name)
}

func kubectl(ctx context.Context, args ...string) error {
	c := cmd.New("/usr/local/bin/kubectl",
		cmd.WithImplicitEnv(),
		cmd.WithArgs(args...))
	return c.Execute(ctx)
}
