.PHONY: up down build logs restart clean

up:
	docker compose up --build -d

down:
	docker compose down -v --remove-orphans

build:
	docker compose build --no-cache

logs:
	docker compose logs -f

restart: down up 