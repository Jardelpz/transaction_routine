up:
	docker-compose up --build

down:
	docker-compose down

reset:
	docker-compose down -v
	docker-compose up --build

test:
	cd .. && go test ./... -short

test-integration:
	cd .. && go test ./... -tags=integration

test-all:
	cd .. && go test ./...
