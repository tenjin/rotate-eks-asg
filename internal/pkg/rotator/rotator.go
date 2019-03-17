package rotator

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/tenjin/rotate-eks-asg/internal/pkg/cmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/exec"
)

const (
	DefaultKubeConfig = "KUBECONFIG=/tmp/.kube/config"
)

func RotateAll(ctx context.Context, cluster string, groups []string) error {
	for _, group := range groups {
		if err := Rotate(ctx, cluster, group); err != nil {
			return err
		}
	}
	return nil
}

func Rotate(ctx context.Context, cluster, asg string) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	a := autoscaling.New(sess)
	group, err := DescribeAutoScalingGroup(a, asg)
	if err != nil {
		return err
	}

	log.Printf("Rotating ASG '%s'...\n", asg)
	for _, id := range group.instanceIds {
		log.Printf("Rotating Instance '%s'...\n", id)
		// TODO
	}
	return nil
}

func RotateOld(ctx context.Context, cluster, asg string) error {
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
	kc, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	nodes, err := kc.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for i, node := range nodes.Items {
		log.Printf("%d: %s", i+1, node.Spec.ProviderID)
	}

	return nil
}
