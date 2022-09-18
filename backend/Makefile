postgres:
	docker run --name postgres14 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root bhagad_house_booking

dropdb:
	docker exec  -it postgres14 dropdb bhagad_house_booking

migrateup:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5433/bhagad_house_booking?sslmode=disable" -verbose up

migratedown:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5433/bhagad_house_booking?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go clean -testcache
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test