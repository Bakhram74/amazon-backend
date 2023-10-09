

docker_run:
	docker run --name amazon -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=123456 -p 5432:5432 -d postgres

createDB:
	docker exec -it amazon createdb --username=postgres --owner=postgres amazon-record

dropDB:
	docker exec -it amazon dropdb amazon-record


migrate:
	migrate create -ext sql -dir db/schema/migration -seq init_schema

migrateUP:
	migrate -path db/schema/migration -database "postgresql://postgres:123456@localhost:5432/amazon-record?sslmode=disable" -verbose up

migrateDown:
	migrate -path db/schema/migration -database "postgresql://postgres:123456@localhost:5432/amazon-record?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run cmd/main.go

mock:
	go generate internal/service/service.go


PHONY:  docker_run createDB dropDB docker_exec migrate migrateUP migrateDown sqlc test server mock
