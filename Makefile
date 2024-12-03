KIND_CLUSTER_NAME := k8sfsm

.PHONY: build
build:
	go build

.PHONY: run
run:
	$(MAKE) build
	./k8sfsm job1/input.yaml job1/output.yaml
	./k8sfsm job2/input.yaml job2/output.yaml
	./k8sfsm job3/input.yaml job3/output.yaml

.PHONY: create-cluster
create-cluster:
	aqua i
	kind create cluster --name $(KIND_CLUSTER_NAME)

.PHONY: delete-cluster
delete-cluster:
	kind delete cluster --name $(KIND_CLUSTER_NAME)
