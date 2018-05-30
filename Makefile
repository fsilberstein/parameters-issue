.PHONY: install
install:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure -v

.PHONY: run
run: 
	@go run main.go

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -o ./app -a -ldflags '-s' -installsuffix cgo main.go

.PHONY: test-unit
test-unit:
	go test -v ./...

.PHONY: test-integration
test-integration:
	go test -v ./... -tags=integration

