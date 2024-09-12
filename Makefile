BINARY_NAME=main.out

init:
	./thegang/get-tools.sh
	./thegang/project-setup.sh

compile:
	go build -o ${BINARY_NAME} main.go

run:
	go build -o ${BINARY_NAME} main.go
	./${BINARY_NAME} server

seed:
	go build -o ${BINARY_NAME} main.go
	./${BINARY_NAME} seed

clean:
	go clean
	rm ${BINARY_NAME}

test:
	go test -cover ./internal/core/controllers ./internal/core/services -coverprofile=coverage.out -v -test.v; go tool cover --html=./coverage.out;

lzl:
	golangci-lint run --fix

lzm:
	go generate ./internal/core/services ./internal/core/vaults

wire:
	cd ./internal/injection && wire

export-table:
	./thegang/manual/export_schema.sh

mig-up:
	go build -o ${BINARY_NAME} main.go
	./${BINARY_NAME} duck up

mig-down:
	go build -o ${BINARY_NAME} main.go
	./${BINARY_NAME} duck down

gen:
	cd openapi && oapi-codegen --config=config.yaml ./spec.yaml