package rotator

import (
	"errors"
	"fmt"
	"log"

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

func DescribeInstanceByInternalDNS(
	ec2Client *ec2.EC2,
	asgClient *autoscaling.AutoScaling,
	instanceInternalDNS string,
) (string, string, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{{
			Name:   aws.String("network-interface.private-dns-name"),
			Values: []*string{aws.String(instanceInternalDNS)},
		}},
	}
	var instanceID string
	err := ec2Client.DescribeInstancesPages(input,
		func(output *ec2.DescribeInstancesOutput, isLast bool) bool {
			instanceID = *(output.Reservations[0].Instances[0].InstanceId)
			return false
		})
	if err != nil {
		return "", "", err
	}
	if instanceID == "" {
		return "", "", errors.New(fmt.Sprintf("%s: No matching instance could be found", instanceInternalDNS))
	}

	log.Printf("Internal DNS '%s' is instance ID '%s'", instanceInternalDNS, instanceID)

	var groupName string
	asgInput := &autoscaling.DescribeAutoScalingInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}
	err = asgClient.DescribeAutoScalingInstancesPages(asgInput,
		func(output *autoscaling.DescribeAutoScalingInstancesOutput, isLast bool) bool {
			for _, instance := range output.AutoScalingInstances {
				groupName = *(instance.AutoScalingGroupName)
				return false
			}
			return true
		})
	if err != nil {
		return "", "", err
	}
	if groupName == "" {
		return "", "", errors.New(fmt.Sprintf("%s: No matching ASG could be found", instanceInternalDNS))
	}

	return instanceID, groupName, nil
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

func DetachInstance(client *autoscaling.AutoScaling, groupId, id string) error {
	log.Printf("Detaching instance '%s' from ASG '%s'...", id, groupId)
	in := &autoscaling.DetachInstancesInput{
		InstanceIds:                    aws.StringSlice([]string{id}),
		AutoScalingGroupName:           aws.String(groupId),
		ShouldDecrementDesiredCapacity: aws.Bool(false),
	}
	_, err := client.DetachInstances(in)
	if err != nil {
		return err
	}
	log.Printf("Instance '%s' detached.", id)
	return nil
}

func TerminateInstanceByID(client *ec2.EC2, id string) error {
	log.Printf("Terminating instance '%s'...", id)
	in := &ec2.TerminateInstancesInput{
		InstanceIds: aws.StringSlice([]string{id}),
	}
	_, err := client.TerminateInstances(in)
	if err != nil {
		return err
	}
	log.Printf("Instance '%s' succesfully terminated.", id)
	return nil
}
