dev:
	go run main.go

docker-build:
	docker login registry.gitlab.com
	docker build -t registry.gitlab.com/pantomath-io/demo-tools .
	docker push registry.gitlab.com/pantomath-io/demo-tools

.PHONY: vendor
vendor:
	rm -rf ./vendor
	govendor init
	govendor add +external