postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root  -e POSTGRES_PASSWORD=daxi1982 -d postgres:alpine
createdb:
	docker exec -it postgres  createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres  dropdb  simple_bank	
migratecreate:
	migrate create -ext sql -dir db/migration -seq int_schema
migrateup:
	migrate -path db/migration -database "postgresql://root:daxi1982@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:daxi1982@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
gitadd:
	git add . 

.PHONY: postgres createdb dropdb migrateup migratedown migratecreate sqlc