run:
	cd app && go run main.go

test: mockgen
	cd app && go test -v ./...

mockgen:
	rm -rf ./app/storage/*_mock.go
	mockgen -source=./app/storage/storage.go -destination=./app/storage/storage_mock.go -package=storage
	rm -rf ./app/util/*_mock.go
	mockgen -source=./app/util/time.go -destination=./app/util/time_mock.go -package=util