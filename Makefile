postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=thobeogalaxy257 -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:thobeogalaxy257@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:thobeogalaxy257@localhost:5432/simple_bank?sslmode=disable" -verbose down
migrateup1ver:
	migrate -path db/migration -database "postgresql://root:thobeogalaxy257@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migratedown1ver:
	migrate -path db/migration -database "postgresql://root:thobeogalaxy257@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/techschool/simplebank/db/sqlc Store
.PHONY: postgres createdb dropdb migrateup migratedown migrateup1ver migratedown1ver sqlc test server mock