.PHONY: seed,compose

seed:
	go run ./cmd/dbseeder

compose:
	docker compose up

compose_build:
	docker compose up --build