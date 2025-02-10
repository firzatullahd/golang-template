run-api:
	go run cmd/api/main.go 

build-api:
	GOARCH=arm64 go build -o golang-api  ./cmd/api/main.go

build-api-dev:
	GOOS=linux GOARCH=amd64 go build -o golang-api  ./cmd/api/main.go


test:
	go test -race -short -coverprofile=./cov.out ./... 
