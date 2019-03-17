# rotate-eks-asg [![Docker Repository on Quay](https://quay.io/repository/tenjin/rotate-eks-asg/status "Docker Repository on Quay")](https://quay.io/repository/tenjin/rotate-eks-asg)

Rolling Cluster Node Upgrades for AWS EKS

**Project Status:** Work in progress and not ready for use at this moment. 

## Use Case

Apply security fixes, rollout new Kubernetes versions, or replace faulty nodes on AWS.

In general terms:

- You run Kubernetes via [AWS EKS](https://aws.amazon.com/eks/)
- Your cluster is made up of [EC2 Auto Scaling Groups (ASG)](https://docs.aws.amazon.com/autoscaling/ec2/userguide/AutoScalingGroup.html)
- You want to replace one or all nodes in those ASGs (e.g. to [activate a new launch configuration](https://docs.aws.amazon.com/autoscaling/ec2/userguide/LaunchConfiguration.html))
- The replacement has to be done gracefully, node-by-node, and respects [availability constraints in your cluster](https://kubernetes.io/docs/tasks/run-application/configure-pdb/)
