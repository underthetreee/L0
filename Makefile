.SILENT:

.PHONY: run
run:
	go run cmd/app/main.go

.PHONY: compose-up
compose-up:
	docker-compose up -d

.PHONY: compose-down
compose-down:
	docker-compose down --remove-orphans

.PHONY: migrate-up
migrate-up: 
	goose -dir migrations postgres ${POSTGRES_URL} up

.PHONY: migrate-down
migrate-down: 
	goose -dir migrations postgres ${POSTGRES_URL} down