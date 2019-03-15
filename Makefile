all: docker
docker:
	docker build -t rotate-eks-asg .
