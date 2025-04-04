NAMESPACE := sales-system
VERSION := 0.0.1
COMPOSE_FILE := infra/compose/docker-compose.yaml

################################################################################
run-dev:
	go run cmd/sales/main.go

build: sales


load:
	hey -m GET -c 100 -n 1000 "http://localhost:8000/v1/test"

metrics:
	expvarmon -ports="localhost:3000" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

help: 
	go run cmd/sales/main.go --help

tidy:
	go mod tidy 
	go mod vendor

key-gen:
	go run cmd/admin/main.go genkey -size=2048

token-gen:
	go run cmd/admin/main.go gentoken -kid=c3550713-13e7-4a53-977a-dd53cbcb7088 -keypath=infra/keys -userid=random_user

################################################################################

sales:
	docker build \
	-f infra/docker/sales.dockerfile \
	-t sales:$(VERSION) \
	--build-arg BUILD_REF=$(VERSION) \
	.


apply:
	kubectl apply -f infra/k8s/base/namespace.yaml
	kubectl apply -f infra/k8s/database/database.yaml 
	kubectl rollout status --namespace=$(NAMESPACE) --watch --timeout=120s sts/database
	kubectl apply -f infra/k8s/sales/sales-deploy.yaml
	

delete:
	kubectl delete -f infra/k8s/sales/sales-deploy.yaml
	kubectl delete -f infra/k8s/database/database.yaml	
	kubectl delete -f infra/k8s/base/namespace.yaml

restart-sales:
	kubectl rollout restart deployment sales-deployment --namespace=$(NAMESPACE) 

restart-database:
	kubectl rollout restart statefulset database --namespace=$(NAMESPACE)

status:
	kubectl get pods --namespace=$(NAMESPACE) --watch -o wide

logs-sales:
	kubectl logs --namespace=$(NAMESPACE) -l app=sales --all-containers=true -f --tail=100 --max-log-requests=6 

logs-database:
	kubectl logs --namespace=$(NAMESPACE) -l app=database --all-containers=true -f --tail=100

describe-sales-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) sales-deployment

describe-sales-pods:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=sales

describe-database:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=database


#-----------------------------------------------------------------------------------------------------------------------
# compose
compose-up:
	docker-compose -f $(COMPOSE_FILE) up -d
compose-down:
	docker-compose -f $(COMPOSE_FILE) down

compose-build: build compose-up
compose-logs:
	docker-compose -f $(COMPOSE_FILE) logs -f