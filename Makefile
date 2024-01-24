postgres:
	docker run --name postgres1 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:latest

createdb:
	docker exec -it postgres1 createdb --username=root --owner=root pet_bank

dropdb:
	docker exec -it postgres1 drop pet_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/pet_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/pet_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go pet-bank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock