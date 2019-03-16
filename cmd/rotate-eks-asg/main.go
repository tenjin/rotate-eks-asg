package main

import (
	"github.com/complex64/go-utils/pkg/ctxutil"
	"github.com/tenjin/rotate-eks-asg/internal/pkg/cmd"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
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
	ctx, _ := ctxutil.ContextWithCancelSignals(os.Kill, os.Interrupt)

	log.Print("Fetching EKS configuration...")
	eksCfgCmd := cmd.New("aws",
		cmd.WithArgs("eks", "update-kubeconfig", "--name", *eksCluster),
		cmd.WithImplicitEnv(),
		cmd.WithEnv(DefaultKubeConfig))
	if err := eksCfgCmd.Execute(ctx); err != nil {
		log.Fatal(err)
	}


}
