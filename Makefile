.PHONY: up down build logs restart clean

up:
	docker compose up --build -d

down:
	docker compose down -v --remove-orphans

build:
	docker compose build --no-cache

logs:
	docker compose logs -f

lint:
	docker run --rm -v $(CURDIR):/app:Z -w /app golangci/golangci-lint:latest golangci-lint run ./...

restart: down up 