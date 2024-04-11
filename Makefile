include .env

db\:migrate\:up:
	@echo "Up migrations"
	@goose -dir ${GOOSE_MIGRATION_DIR} ${GOOSE_DRIVER} ${GOOSE_DBSTRING} up

db\:migrate\:down:
	@echo "Down migrations"
	@goose -dir ${GOOSE_MIGRATION_DIR} ${GOOSE_DRIVER} ${GOOSE_DBSTRING} down

db\:seed\:up:
	@echo "Up seeds"
	@goose -no-versioning -dir ${GOOSE_SEED_DIR} ${GOOSE_DRIVER} ${GOOSE_DBSTRING} up

db\:reset:
	@echo "Reset database"
	@goose -dir ${GOOSE_MIGRATION_DIR} ${GOOSE_DRIVER} ${GOOSE_DBSTRING} reset
	@make db:migrate:up && make db:seed:up
	@sqlc generate
