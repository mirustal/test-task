migrationUp:
	goose -dir db/migrations/postgres create first_migration sql

serviceUp:
	docker-compose down
	docker-compose up -d

serviceRestart:
	docker-compose up --build --no-deps bank_service
