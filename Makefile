all: test
docker:
	docker build -t rotate-eks-asg .
deps:
	GO111MODULE=on go build ./cmd/...
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor
install:
	GO111MODULE=on go install -mod=vendor ./cmd/rotate-eks-asg
	GO111MODULE=on go install -mod=vendor ./cmd/rotate-eks-instance
