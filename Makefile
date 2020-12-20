MIGRATION_ARG ?= up

image:
	docker build . --target dev -t auth-backend:api-dev

build:
	docker-compose run --rm api go build -i -o api

deps:
	docker-compose run --rm api go mod tidy
	docker-compose run --rm api go mod vendor

setup: image deps build
	docker-compose run --rm auth_db mysql -uroot -psecret -hauth_db -e "CREATE DATABASE IF NOT EXISTS auth"
	docker-compose run --rm api ./api migrate install
	docker-compose run --rm api ./api migrate

local: local-down build
	docker-compose up api

local-down:
	docker-compose rm -sf

test:
	docker-compose run --rm api go test -cover -v ./... -tags nojira

migrate: build
	docker-compose run --rm api ./api migrate $(MIGRATION_ARG)

migration: build
	docker-compose run --rm api ./api make:migration $(name)

seed: build
	docker-compose run --rm api ./api seed $(SEED_ARG)

db-browse:
	docker-compose exec auth_db mysql -uroot -psecret

push:
	aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin 996760811179.dkr.ecr.us-east-2.amazonaws.com
	docker build . --target final -t auth-backend
	docker tag auth-backend:latest 996760811179.dkr.ecr.us-east-2.amazonaws.com/auth-backend:latest
	docker push 996760811179.dkr.ecr.us-east-2.amazonaws.com/auth-backend:latest
