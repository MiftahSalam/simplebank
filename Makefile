postgres:
	docker run --name postgres_simplebank -p 5243:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it postgres_simplebank createdb --username=root --owner=root simplebank

dropdb:
	docker exec -it postgres_simplebank dropdb simplebank

migrateup:
	migrate --path db/migrations --database "postgresql://root:secret@localhost:5243/simplebank?sslmode=disable" --verbose up

migratedown:
	migrate --path db/migrations --database "postgresql://root:secret@localhost:5243/simplebank?sslmode=disable" --verbose down

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown test
