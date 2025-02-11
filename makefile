VERSION := 0.0.1
################################################################################
run-dev:
	go run cmd/sales/main.go

build: sales

sales:
	docker build \
	-f infra/docker/sales.dockerfile \
	-t sales:$(VERSION) \
	--build-arg BUILD_REF=$(VERSION) \
	.