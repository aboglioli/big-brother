APP=main

build: lint clean
	go build -o ${APP} cmd/main.go

run:
	go run -race cmd/main.go

test:
	go test ./...

short-test:
	go test -short ./...

lint:
	gofmt -s -w .
	go vet ./...

clean:
	go clean -testcache ./...