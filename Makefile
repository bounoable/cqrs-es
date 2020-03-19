default: services docker-test

test:
	go test ./...

docker-test:
	docker-compose run --rm library make test

services:
	docker-compose up -d nats nats_streaming

stop:
	docker-compose down

.PHONY: test docker-test services
