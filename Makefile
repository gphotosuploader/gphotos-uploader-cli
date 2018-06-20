start:
	go run main.go

.PHONY: vendor
vendor:
	rm -rf ./vendor
	govendor init
	govendor add +external