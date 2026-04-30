DB_URL=postgresql://root:root@localhost:5433/postgres?sslmode=disable
postgres:
	docker run --name my_postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -e POSTGRES_DB=postgres -p 5433:5432 -d postgres:alpine

createdb:
	docker exec -it my_postgres createdb --username=root --owner=root postgres

dropdb:
	docker exec -it my_postgres dropdb --username=root --owner=root postgres

migrateup:
	migrate -path db/migrations/ -database "$(DB_URL)" -verbose up 

migratedown:
	migrate -path db/migrations/ -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migratedown migrateup sqlc