# rotate-eks-asg [![Docker Repository on Quay](https://quay.io/repository/tenjin/rotate-eks-asg/status "Docker Repository on Quay")](https://quay.io/repository/tenjin/rotate-eks-asg)

Rolling Cluster Node Upgrades for AWS EKS

**Project Status:** Used in production at Tenjin, some caveats apply.

## Use Case

Apply security fixes, rollout new Kubernetes versions, or replace faulty nodes on AWS.

In general terms:

- You run Kubernetes via [AWS EKS](https://aws.amazon.com/eks/)
- Your cluster is made up of [EC2 Auto Scaling Groups (ASG)](https://docs.aws.amazon.com/autoscaling/ec2/userguide/AutoScalingGroup.html)
- You want to replace one or all nodes in those ASGs (e.g. to [activate a new launch configuration](https://docs.aws.amazon.com/autoscaling/ec2/userguide/LaunchConfiguration.html))
- The replacement has to be done gracefully, node-by-node, and respects [availability constraints in your cluster](https://kubernetes.io/docs/tasks/run-application/configure-pdb/)

## Usage

You can run this tool from your CI or locally. Typically we bundle it as a script and inject secrets within the CI.

Example using standard AWS SDK credentials and an assumed role:

```bash
#!/bin/bash
set -ex
docker run --rm -it \
    -e ACCESS_KEY_ID=${ACCESS_KEY_ID:?}
    -e SECRET_ACCESS_KEY=${SECRET_ACCESS_KEY:?}
    -e ROLE_ARN=${ROLE_ARN:?}
    -e CLUSTER=your-cluster-name \
    -e AUTOSCALING_GROUPS=${AUTOSCALING_GROUP:?} \
    rotate-eks-asg:latest
```
