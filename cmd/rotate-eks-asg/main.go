package main

import (
	"log"
	"os"

	"github.com/complex64/go-utils/pkg/ctxutil"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/tenjin/rotate-eks-asg/internal/pkg/rotator"
)

var (
	cluster = kingpin.Arg("asg", "EKS cluster formed by the Auto Scaling Groups (ASG)").Required().String()
	asgs    = kingpin.Arg("cluster", "EKS Auto Scaling Groups to rotate").Required().Strings()
)

func main() {
	kingpin.Parse()
	ctx, cancel := ctxutil.ContextWithCancelSignals(os.Kill, os.Interrupt)
	defer cancel()
	for _, group := range *asgs {
		if err := rotator.Rotate(ctx, *cluster, group); err != nil {
			log.Fatal(err)
		}
	}
}
