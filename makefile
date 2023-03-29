run:
	go run cmd/app/main.go

mock:
	mockery --all

test:
	go test ./... --cover

lint:
	golangci-lint run

performance:
	go run cmd/performance/performance.go