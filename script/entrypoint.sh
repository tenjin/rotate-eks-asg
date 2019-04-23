#!/bin/sh
set -ex

# Render a AWS credentials with user and role for `awsudo`.
DEFAULT_REGION=us-east-1
REGION=${REGION:-${DEFAULT_REGION}}

if [[ ! -f ~/.aws/credentials ]]; then
    mkdir -p ~/.aws
    cat <<EOF >~/.aws/credentials
[user]
region = ${REGION}
aws_access_key_id = ${ACCESS_KEY_ID:?}
aws_secret_access_key = ${SECRET_ACCESS_KEY:?}
[role]
region = ${REGION}
source_profile = user
role_arn = ${ROLE_ARN:?}
EOF
fi

export KUBECONFIG=/tmp/.kube/config
awsudo -u role aws eks update-kubeconfig --name ${CLUSTER:?}
awsudo -u role rotate-eks-asg ${AUTOSCALING_GROUPS:?}
