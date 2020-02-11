APP=main
CONTAINERS=postgres postgres-pgadmin redis redis-commander
POSTGRES_CONTAINER=postgres
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USERNAME=admin

build: lint clean
	go build -o ${APP} cmd/main.go

run:
	go run -race cmd/main.go

test:
	# go test ./...
	gotestsum

vtest:
	go test ./... -v

stest:
	# go test -short ./...
	gotestsum ./... -short

lint:
	gofmt -s -w .
	go vet ./...

clean:
	go clean -testcache ./...

# Environment
up:
	docker-compose up -d ${CONTAINERS}

down:
	docker-compose down -v --remove-orphans

create-databases:
	docker-compose exec \
		${POSTGRES_CONTAINER} \
		psql -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USERNAME} \
		-f migrations/databases.sql

migrate:
	go run cmd/migrate/main.go
