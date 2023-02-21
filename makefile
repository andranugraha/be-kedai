run:
	go run cmd/app/main.go

mock:
	mockery --all

test:
	go test ./... --cover