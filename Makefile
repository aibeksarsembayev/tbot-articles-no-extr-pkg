build:
	docker-compose up -d
	docker image prune

build-dev:
	docker-compose -f docker-compose-dev.yml up
	docker image prune
	go run . 

build-dbtest:
	docker-compose -f docker-compose-dbtest.yml up
	docker image prune

stop:
	docker-compose down

delete:
	docker rmi alpine:latest sber-invest-bot_app:latest postgres:latest golang:1.19 

test:
	go test ./... -v -cover