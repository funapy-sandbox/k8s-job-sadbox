
CLUSTER_NAME ?= funapy-sandbox/k8s-job-sadbox

.PHONY: k3d/start
k3d/start:
	k3d cluster create $(CLUSTER_NAME)

.PHONY: k3d/stop
k3d/stop:
	k3d cluster delete $(CLUSTER_NAME)
