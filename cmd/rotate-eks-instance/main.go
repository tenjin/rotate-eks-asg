package main

import (
	"log"
	"os"

	"github.com/complex64/go-utils/pkg/ctxutil"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/tenjin/rotate-eks-asg/internal/pkg/rotator"
)

var (
	name      = kingpin.Arg("name", "Internal DNS of EKS instance to rotate").Required().String()
	removeNode = kingpin.Flag("remove", "Remove instance, don't provision a replacement").Default("false").Bool()
)

func init() {
	_ = os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
}

func main() {
	kingpin.Parse()
	ctx, cancel := ctxutil.ContextWithCancelSignals(os.Kill, os.Interrupt)
	defer cancel()
	if err := rotator.RotateByInternalDNS(ctx, *name, *removeNode); err != nil {
		log.Fatal(err)
	}
}
