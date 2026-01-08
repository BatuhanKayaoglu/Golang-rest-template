setup:
	go get -u github.com/swaggo/swag/cmd/swag
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g ./cmd/server/main.go -o ./docs
	go get -u github.com/swaggo/gin-swagger
	go get -u github.com/swaggo/files

build-docker:
	docker compose build --no-cache

run-local:
	docker compose up db redis mongo -d
	REDIS_HOST=localhost \
	POSTGRES_DB=go_app_dev \
	POSTGRES_USER=docker \
	POSTGRES_PASSWORD=password \
	POSTGRES_PORT=5435 \
	JWT_SECRET_KEY=ObL89O3nOSSEj6tbdHako0cXtPErzBUfq8l8o/3KD9g=INSECURE \
	API_SECRET_KEY=cJGZ8L1sDcPezjOy1zacPJZxzZxrPObm2Ggs1U0V+fE=INSECURE \
	POSTGRES_HOST=localhost \
	MONGO_HOST=localhost \
	go run cmd/server/main.go

up:
	docker compose up

down:
	docker compose down

restart:
	docker compose restart

build:
	go build -v ./...

test:
	go test -v ./... -race -cover

clean:
	docker stop go-rest-api-template
	docker stop dockerPostgres
	docker rm go-rest-api-template
	docker rm dockerPostgres
	docker rm dockerRedis
	docker image rm golang-rest-api-template-backend
	rm -rf .dbdata

rebuild:
	docker compose down && docker compose build backend && docker compose up

