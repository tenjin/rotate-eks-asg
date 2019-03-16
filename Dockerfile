FROM golang:alpine as buildenv
WORKDIR /go/src/github.com/tenjin/rotate-eks-asg
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /rotate-eks-asg ./cmd/rotate-eks-asg

# Extract `rotate-eks-asg` and ship with rest of tooling:

FROM python:3-alpine
RUN apk add --no-cache \
    curl \
    git

ARG AWSCLI_VERSION=1.16.58
ARG KUBECTL_VERSION=1.13.0

# https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html
ARG AWSAUTHENTICATOR_URL=https://amazon-eks.s3-us-west-2.amazonaws.com/1.11.5/2018-12-06/bin/linux/amd64/aws-iam-authenticator

# https://kubernetes.io/docs/tasks/tools/install-kubectl
ARG KUBECTL_URL=https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl

# `awsudo` and `aws`
RUN pip install -q \
      awscli==${AWSCLI_VERSION} \
      git+https://github.com/tenjin/awsudo.git

ADD ${KUBECTL_URL} /usr/local/bin/kubectl
ADD ${AWSAUTHENTICATOR_URL} /usr/local/bin/aws-iam-authenticator

COPY --from=buildenv /rotate-eks-asg /usr/local/bin/

RUN chmod +x \
    /usr/local/bin/kubectl \
    /usr/local/bin/aws-iam-authenticator

ENTRYPOINT ["/usr/local/bin/rotate-eks-asg"]
