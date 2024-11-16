migrationUp:
	goose -dir db/migrations/postgres create first_migration sql

serviceUp:
	docker-compose down
	docker-compose up --build