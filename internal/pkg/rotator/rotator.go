package rotator

import (
	"context"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/tenjin/rotate-eks-asg/internal/pkg/cmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/exec"
)

const (
	DefaultKubeConfig = "KUBECONFIG=/tmp/.kube/config"
)

func Rotate(ctx context.Context, cluster, asg string) error {
	log.Print("Fetching EKS configuration...")
	eksCfgCmd := cmd.New("aws",
		cmd.WithArgs("eks", "update-kubeconfig", "--name", cluster),
		cmd.WithImplicitEnv(),
		cmd.WithEnv(DefaultKubeConfig))
	if err := eksCfgCmd.Execute(ctx); err != nil {
		log.Fatal(err)
	}

	config, err := clientcmd.BuildConfigFromFlags("", "/tmp/.kube/config") // TODO
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for i, node := range nodes.Items {
		log.Printf("%d: %s", i+1, node.Name)
	}

	return nil
}
