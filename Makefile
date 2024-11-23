DB_URL=postgres://postgres:Abdu0811@localhost:5432/postgres?sslmode=disable


migrate:
	migrate create -dir migrations -ext sql db
migrateup:
	migrate -path migrations -database ${DB_URL} up
migratedown:
	migrate -path migrations -database ${DB_URL} down
migrateforce:
	migrate -path migrations -database ${DB_URL} force

run_db:
	sudo docker compose up -d db

run:
	sudo docker compose up --build

stop:
	sudo docker-compose down

clean:
	sudo docker-compose down -v
swag:
	swag init -g internal/http/handler/auth.go