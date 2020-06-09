package rotator

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/exec"
)

func RotateAll(ctx context.Context, groups []string) error {
	for _, group := range groups {
		if err := Rotate(ctx, group); err != nil {
			return err
		}
	}
	return nil
}

func Rotate(ctx context.Context, groupId string) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	asgClient := autoscaling.New(sess)
	ec2Client := ec2.New(sess)
	group, err := DescribeAutoScalingGroup(asgClient, groupId)
	if err != nil {
		return err
	}
	k8s, err := NewKubernetesClient()
	if err != nil {
		return err
	}

	log.Printf("Rotating ASG '%s'...\n", groupId)
	for _, id := range group.instanceIds {
		log.Printf("Rotating Instance '%s'...\n", id)
		if err := RotateInstance(ctx, k8s, asgClient, ec2Client, groupId, id, false); err != nil {
			return err
		}
	}
	return nil
}

func RotateByInternalDNS(ctx context.Context, instanceInternalIP string, removeNode bool) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	asgClient := autoscaling.New(sess)
	ec2Client := ec2.New(sess)
	instanceID, groupID, err := DescribeInstanceByInternalDNS(ec2Client, asgClient, instanceInternalIP)
	if err != nil {
		return err
	}
	k8s, err := NewKubernetesClient()
	if err != nil {
		return err
	}
	return RotateInstance(ctx, k8s, asgClient, ec2Client, groupID, instanceID, removeNode)
}

func RotateInstance(
	ctx context.Context,
	k8s *kubernetes.Clientset,
	asg *autoscaling.AutoScaling,
	ec2 *ec2.EC2,
	groupId string,
	instanceId string,
	removeNode bool,
) error {
	name, err := GetNodeNameByInstanceID(ctx, k8s, instanceId)
	if err != nil {
		return err
	}
	if err := CordonNodeByName(ctx, name); err != nil {
		return err
	}
	nodeSet, err := GetClusterNodeSet(ctx, k8s)
	if err != nil {
		return err
	}
	if err := DetachInstance(asg, groupId, instanceId, removeNode); err != nil {
		return err
	}

	if !removeNode {
		if err := AwaitNewNodeReady(ctx, k8s, nodeSet); err != nil {
			return err
		}
	}

	if err := DrainNodeByName(ctx, name); err != nil {
		return err
	}
	if err := TerminateInstanceByID(ec2, instanceId); err != nil {
		return err
	}
	return nil
}
