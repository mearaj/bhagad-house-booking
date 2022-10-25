DB_URL=postgresql://root:secret@localhost:5432/bhagad_house_booking?sslmode=disable

postgres:
	docker run --name postgres14 --network bhagad_house_network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root bhagad_house_booking

dropdb:
	docker exec  -it postgres14 dropdb bhagad_house_booking

migrateup:
	migrate -path ./db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path ./db/migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

mock:
	mockgen -package mock -destination db/mock/store.go github.com/mearaj/bhagad-house-booking/common/db/sqlc Store

test:
	go clean -testcache
	go test -v -cover ./...

.PHONY: postgres dropdb createdb migrateup migratedown sqlc test