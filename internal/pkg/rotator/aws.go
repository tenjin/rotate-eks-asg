package rotator

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Group struct {
	launchConfigurationName string
	instanceIds             []string
}

func DescribeAutoScalingGroup(client *autoscaling.AutoScaling, name string) (*Group, error) {
	group, err := getAutoScalingGroup(client, name)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(group.Instances))
	for _, i := range group.Instances {
		ids = append(ids, *i.InstanceId)
	}
	g := &Group{
		launchConfigurationName: *group.LaunchConfigurationName,
		instanceIds:             ids,
	}
	return g, nil
}

func getAutoScalingGroup(client *autoscaling.AutoScaling, name string) (*autoscaling.Group, error) {
	in := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: aws.StringSlice([]string{name}),
		MaxRecords:            aws.Int64(1),
	}
	out, err := client.DescribeAutoScalingGroups(in)
	if err != nil {
		return nil, err
	}
	if len(out.AutoScalingGroups) != 1 {
		return nil, fmt.Errorf("expected exactly 1 ASG description for '%s' got %d", name, len(out.AutoScalingGroups))
	}
	return out.AutoScalingGroups[0], nil
}

func DetachInstance(client *autoscaling.AutoScaling, group, id string) error {
	in := &autoscaling.DetachInstancesInput{
		InstanceIds:                    aws.StringSlice([]string{id}),
		AutoScalingGroupName:           aws.String(group),
		ShouldDecrementDesiredCapacity: aws.Bool(false),
	}
	_, err := client.DetachInstances(in)
	if err != nil {
		return err
	}
	return nil
}

func TerminateInstanceByID(client *ec2.EC2, id string) error {
	in := &ec2.TerminateInstancesInput{
		InstanceIds: aws.StringSlice([]string{id}),
	}
	_, err := client.TerminateInstances(in)
	if err != nil {
		return err
	}
	return nil
}
