postgres:
	docker run --name postgres_simplebank -p 5243:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it postgres_simplebank createdb --username=root --owner=root simplebank

dropdb:
	docker exec -it postgres_simplebank dropdb simplebank

create-migrate:
	migrate create -ext sql -dir db/migrations -seq $(name)

migrateup:
	migrate --path db/migrations --database "postgresql://root:secret@localhost:5243/simplebank?sslmode=disable" --verbose up

migratedown:
	migrate --path db/migrations --database "postgresql://root:secret@localhost:5243/simplebank?sslmode=disable" --verbose down

migrateup1:
	migrate --path db/migrations --database "postgresql://root:secret@localhost:5243/simplebank?sslmode=disable" --verbose up 1

migratedown1:
	migrate --path db/migrations --database "postgresql://root:secret@localhost:5243/simplebank?sslmode=disable" --verbose down 1

test:
	go test -v -cover ./...

mockdb:
	mockgen -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store

server-dev:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 test
