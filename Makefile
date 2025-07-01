Container:
	docker run --name jansi -p 5433:5432 -e POSTGRES_PASSWORD=SituBen -e POSTGRES_USER=root -d postgres

Createdb:
	docker exec -it nilesh createdb --username=root --owner=root Hiten

Dropdb:
	docker exec -it nilesh dropdb Hiten

MigrateUp:
	migrate -path db/migration -database "postgres://root:SituBen@localhost:5433/Hiten?sslmode=disable" -verbose up

MigrateDown:
	migrate -path db/migration -database "postgres://root:SituBen@localhost:5433/Hiten?sslmode=disable" -verbose down

Sqlc:
	sqlc generate

Test:
	go test -v -cover ./...

.PHONY:	Container	Createdb	Dropdb	MigrateDown	MigrateUp	Sqlc	Test