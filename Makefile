APP_NAME=warehouse

test-integration:
	go test ./... -v -tags=integration