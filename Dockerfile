FROM golang:alpine as buildenv
WORKDIR /go/src/github.com/tenjin/rotate-eks-asg
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /rotate-eks-asg ./cmd/rotate-eks-asg

# Extract `rotate-eks-asg` and ship with rest of tooling:

FROM python:3-alpine
RUN apk add --no-cache git
RUN pip install -q git+https://github.com/tenjin/awsudo.git

# https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html
ARG AUTHENTICATOR=https://amazon-eks.s3-us-west-2.amazonaws.com/1.11.5/2018-12-06/bin/linux/amd64/aws-iam-authenticator
ADD ${AUTHENTICATOR} /usr/local/bin/aws-iam-authenticator
RUN chmod +x /usr/local/bin/aws-iam-authenticator

COPY --from=buildenv /rotate-eks-asg /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/rotate-eks-asg"]
