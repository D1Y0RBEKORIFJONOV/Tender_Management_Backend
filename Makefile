DB_URL=postgres://postgres:+_+diyor2005+_+@localhost:5432/udevs?sslmode=disable


migrate:
	migrate create -dir migrations -ext sql db
migrateup:
	migrate -path migrations -database ${DB_URL} up
migratedown:
	migrate -path migrations -database ${DB_URL} down
migrateforce:
	migrate -path migrations -database ${DB_URL} force 1
