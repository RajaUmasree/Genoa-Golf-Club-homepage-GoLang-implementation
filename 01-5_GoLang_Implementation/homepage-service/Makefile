.PHONY: build run test vet fmt tidy docker-build

build:
	go build -o bin/homepage-service ./cmd/server

run: build
	./bin/homepage-service

test:
	go test ./... -v -cover

vet:
	go vet ./...

fmt:
	gofmt -l -w .

tidy:
	go mod tidy

docker-build:
	docker build -t genoagolfclub/homepage-service:latest .
