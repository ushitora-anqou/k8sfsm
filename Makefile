KIND_CLUSTER_NAME := k8sfsm

.PHONY: run
run:
	go run main.go

.PHONY: create-cluster
create-cluster:
	aqua i
	kind create cluster --name $(KIND_CLUSTER_NAME)

.PHONY: delete-cluster
delete-cluster:
	kind delete cluster --name $(KIND_CLUSTER_NAME)
