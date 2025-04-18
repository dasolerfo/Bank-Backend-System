postgres: 
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=Songoku13 -d postgres:12-alpine
startDB:
	docker start postgres12

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:Songoku13@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:Songoku13@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:Songoku13@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:Songoku13@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

upgradesqlc:
	brew upgrade sqlc

test: 
	go test -v -cover ./...

mockgen:
	mockgen -package mockdb -destination db/mock/store.go ./db/model Store

serve: 
	go run main.go
	
	
.PHONY: createdb startDB dropdb postgres migrateup migratedown sqlc test serve upgradesqlc


