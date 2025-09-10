run:
	cd app && go run main.go

lint: mockgen
	cd app && golangci-lint run

test: mockgen
	cd app && go test -v ./...

mockgen:
	rm -rf ./app/storage/*_mock.go
	mockgen -source=./app/storage/storage.go -destination=./app/storage/storage_mock.go -package=storage