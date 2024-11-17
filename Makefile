DB_URL=postgres://postgres:2005@localhost:5432/udevs?sslmode=disable


migrate:
	migrate create -dir migrations -ext sql db
migrateup:
	migrate -path migrations -database ${DB_URL} up
migratedown:
	migrate -path migrations -database ${DB_URL} down
migrateforce:
	migrate -path migrations -database ${DB_URL} force

run_db:
	docker-compose up -d db

run:
	docker-compose up --build

stop:
	docker-compose down

clean:
	docker-compose down -v