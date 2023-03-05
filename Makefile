postgres:
	docker run --name  postgres14  -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=assim -d postgres:14-alpine

createdb:
	winpty docker exec -it postgres14 createdb --username=root --owner=root simple_bank

dropdb:
	winpty docker exec -it postgres14 dropdb simple_bank

migrateup:
	migrate -path db/migrations  -database "postgresql://root:assim@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:assim@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migrations  -database "postgresql://root:assim@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migrations -database "postgresql://root:assim@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc: 
	docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb  -destination db/mock/store.go github.com/prepStation/simple_bank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc server mock migrateup1 migratedown1