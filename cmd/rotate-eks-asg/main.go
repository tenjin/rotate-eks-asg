package main

import (
	"log"
	"os"

	"github.com/complex64/go-utils/pkg/ctxutil"
	"gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/tenjin/rotate-eks-asg/internal/pkg/cmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth/exec"
)

const (
	DefaultKubeConfig = "KUBECONFIG=/tmp/.kube/config"
)

var (
	eksCluster = kingpin.Arg("asg", "EKS cluster the ASGs belong to").Required().String()
	asgNames   = kingpin.Arg("cluster", "Names of the EKS ASGs to rotate").Required().Strings()
)

func main() {
	kingpin.Parse()
	ctx, cancel := ctxutil.ContextWithCancelSignals(os.Kill, os.Interrupt)
	defer cancel()

	log.Print("Fetching EKS configuration...")
	eksCfgCmd := cmd.New("aws",
		cmd.WithArgs("eks", "update-kubeconfig", "--name", *eksCluster),
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

	for i,node := range nodes.Items {
		log.Printf("%d: %s", (i+1), node.Name)
	}
}
