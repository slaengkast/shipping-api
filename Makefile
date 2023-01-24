BINARY=shipping-api-server
VERSION := $(shell git rev-parse --short HEAD)

.PHONY: build
build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY} ./cmd/...

.PHONY: test
test:
	GOARCH=amd64 GOOS=linux go test -parallel=8 -coverprofile=coverage.out ./...

run: build
	./${BINARY} --port 8080 --logLevel debug

fmt:
	go fmt ./...

docker:
	docker build -t shipping-api:${VERSION} .

clean:
	go clean
	rm ${BINARY}

get_version:
	@echo ${VERSION}
