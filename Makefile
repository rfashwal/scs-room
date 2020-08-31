## Run the unit-tests only.
test:
	go test ./... -short

lint:
	golangci-lint run --fast

## Run the unit-tests with race detection on.
race:
	go test -race ./... -short

image:
	docker build -t rfashwal/scs-rooms -f Dockerfile .
