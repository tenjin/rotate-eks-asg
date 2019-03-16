all: test
docker:
	docker build -t rotate-eks-asg .
deps:
	GO111MODULE=on go build ./cmd/...
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor
install:
	GO111MODULE=on go install -mod=vendor ./cmd/pinger
test:
	go install ./cmd/rotate-eks-asg
	awsudo -u staging rotate-eks-asg staging-us asg-name
