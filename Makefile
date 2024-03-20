compose:
	@docker-compose build && docker-compose up --force-recreate

test:
	@docker run --rm -v $(PWD):/app -w /app golang:1.22.1-alpine go test ./...

lint:
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.57.1 golangci-lint run
