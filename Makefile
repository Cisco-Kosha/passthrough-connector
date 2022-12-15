BINARY_NAME=passthrough-connector

hello:
	echo "hello"

build:
	 GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin main.go
	 GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/${BINARY_NAME}-linux-arm main.go
	GOOS=linux GOARCH=arm64 go build -o bin/${BINARY_NAME}-linux-arm64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/${BINARY_NAME}-freebsd-386 main.go


clean:
	echo "execute this command to clean up the binaries generated"
	 go clean
	 rm ${BINARY_NAME}-darwin
	 rm ${BINARY_NAME}-linux

test:
	 go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

all: hello build

