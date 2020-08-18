check_install:
	which swagger || GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
	GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models

run_swagger: swagger
	go run main.go

build_and_push:
	docker-compose -f docker/docker-compose.yml build
	docker-compose -f docker/docker-compose.yml push

run:
	docker-compose -f docker/docker-compose.yml up