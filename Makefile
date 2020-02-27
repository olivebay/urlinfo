SHELL := /bin/bash
.DEFAULT_GOAL := help

##@ HELPERS
help:  ## display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ ALL-IN-ONES
all: k3s_install deploy_app ## Spin up a K3s cluster and install urlinfo application

destroy: destroy_cluster ## Destroy cluster

test: test ## Test urlinfo app

loadstart: loadstart ## Start load testing

loadstop: loadstop ## Stop load testing

## K3S
k3s_install:
		@echo ""
		@echo "Installing K3s.."
		@echo ""
		curl -sfL https://get.k3s.io | K3S_KUBECONFIG_MODE="644" sh -s - && echo "" && echo "Sleeping for 60s to allow K3s infrastructure Pods to come up.." && echo "" && sleep 60
		@echo ""

destroy_cluster:
		@echo "Destroying K3s cluster.."
		@echo ""
		/usr/local/bin/k3s-uninstall.sh

deploy_app:
		@echo "Deploying urlinfo application.."
		@echo ""
		kubectl apply -f https://raw.githubusercontent.com/olivebay/urlinfo/master/kubernetes/deployment.yml
		watch -n1 "kubectl get pods -A -o wide"

loadstart:
		@echo "Starting load test.."
		for i in {1..6}; do kubectl run --restart=Never loadtest$$i --image=appropriate/curl -- sh -c "while true; do curl -X GET -i http://apisvc/urlinfo/1/domain.com; done"; done

loadstop:
		@echo "Stopping loadtest.."
		for i in {1..6}; do kubectl delete pod --grace-period=0 --force loadtest$$i; done
test:
		@echo "Testing functionality.."
		curl -X GET -i http://127.0.0.1:31000/urlinfo/1/domain.com