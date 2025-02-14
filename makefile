run:
	go run cmd/main.go

lint:
	golangci-lint run -c .golangci.yml

unit-test:
	go test ./... -short

generate-mocks:
	mockgen -source=./internal/observer/observer.go -destination=./internal/observer/mock_observer.go -package=observer
	mockgen -source=./internal/database/database.go -destination=./internal/database/mock_database.go -package=database
