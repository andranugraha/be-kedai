run:
	go run cmd/app/main.go

mock:
	mockery --all

test:
	go test ./... --cover

lint:
	golangci-lint run

performance-report:
	go run cmd/performance/performance.go

profile:
# DEV
	go tool pprof -http=:8081 http://localhost:8080/debug/pprof/profile
# PROD
# go tool pprof -http=:8081 https://dev-kedai-y3gq8.ondigitalocean.app/debug/pprof/profile