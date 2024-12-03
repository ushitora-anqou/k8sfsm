KIND_CLUSTER_NAME := k8sfsm

.PHONY: build
build:
	go build

.PHONY: run
run:
	$(MAKE) build
	./k8sfsm examples/job1/input.yaml examples/job1/output.yaml
	./k8sfsm examples/job2/input.yaml examples/job2/output.yaml
	./k8sfsm examples/job3/input.yaml examples/job3/output.yaml

.PHONY: create-cluster
create-cluster:
	aqua i
	kind create cluster --name $(KIND_CLUSTER_NAME)

.PHONY: delete-cluster
delete-cluster:
	kind delete cluster --name $(KIND_CLUSTER_NAME)
