ENVIRONMENT := dev
VERSION := 0.1.0
TAG := ${ENVIRONMENT}-${VERSION}
BUILD_INFO := "dapr checkin demo"
KV_NAME=kv-xrgol23o5l5pw
ACR_LOGIN_NAME := acrxrgol23o5l5pw
ACR_URI := acrxrgol23o5l5pw.azurecr.io
FRONTEND_SERVICE_PORT := 8000
BACKEND_SERVICE_PORT := 8001

#### export shell env vars #####
include .env

build_frontend:
	docker build -t ${ACR_URI}/frontend:${TAG} --build-arg SERVICE_NAME="frontend" --build-arg SERVICE_PORT=${FRONTEND_SERVICE_PORT} --build-arg BUILD_INFO=${BUILD_INFO} --build-arg VERSION=${VERSION} -f Dockerfile .
	docker image prune -f

build_backend:
	docker build -t ${ACR_URI}/backend:${TAG} --build-arg SERVICE_NAME="backend" --build-arg SERVICE_PORT=${BACKEND_SERVICE_PORT} --build-arg BUILD_INFO=${BUILD_INFO} --build-arg VERSION=${VERSION} -f Dockerfile .
	docker image prune -f

build: 
	make build_frontend
	make build_backend

push_frontend:
	docker login ${ACR_URI} -u ${ACR_LOGIN_NAME} -p ${acrAdminPassword}
	docker push ${ACR_URI}/frontend:${TAG}

push_backend:
	docker login ${ACR_URI} -u ${ACR_LOGIN_NAME} -p ${acrAdminPassword}
	docker push ${ACR_URI}/backend:${TAG}

push:
	make push_frontend 
	make push_backend

deploy_frontend:
	@if [ -z $(kubectl get deployment frontend) = *"Error from server (NotFound)"* ]; then\
		kubectl apply -f ./manifests/deploy.frontend.yml;\
	fi
	@if [ -z $(kubectl get deployment frontend) != *"Error from server (NotFound)"* ]; then\
		kubectl delete deploy frontend;\
		kubectl apply -f ./manifests/deploy.frontend.yml;\
	fi

deploy_backend:
		@if [ -z $(kubectl get deployment backend) = *"Error from server (NotFound)"* ]; then\
		kubectl apply -f ./manifests/deploy.backend.yml;\
	fi
	@if [ -z $(kubectl get deployment backend) != *"Error from server (NotFound)"* ]; then\
		kubectl delete deploy backend;\
		kubectl apply -f ./manifests/deploy.backend.yml;\
	fi

deploy: 
	make deploy_frontend 
	make deploy_backend

logs_frontend:
	kubectl logs $$(kubectl get pods --selector=app=frontend --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}') frontend -f

logs_backend:
	kubectl logs $$(kubectl get pods --selector=app=backend --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}') backend -f

deploy_dapr_components:
	# deploy AAD pod identity 
	kubectl apply -f ./manifests/azure.pod.identity.yml
	kubectl apply -f ./components/secretstore.yml
	
	# apply ai forwarder manifests
	kubectl apply -f ./manifests/appinsights.forwarder.yml
	# kubectl apply -f ./manifests/exporter.appinsights.yml
	kubectl apply -f ./components/tracing.yml

	# apply keda backend scaler
	kubectl apply -f ./manifests/keda.backend.scaler.yml

	# apply dapr components
	# kubectl apply -f ./components/statestore.yml
	kubectl apply -f ./components/servicebus.queue.binding.yml
	kubectl apply -f ./components/cosmosdb.store.binding.yml

all:
	make build
	make push
	make deploy
