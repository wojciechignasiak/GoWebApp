build: build-macos build-docker-amd64 build-docker-arm64

build-macos:
	GOARCH=arm64 GOOS=darwin go build -o bin/app-macos cmd/application/main.go

build-docker-amd64:
	GOARCH=amd64 GOOS=linux go build -o bin/app-linux-amd64 cmd/application/main.go

build-docker-arm64:
	GOARCH=arm64 GOOS=linux go build -o bin/app-linux-arm64 cmd/application/main.go

runmac: build
	./bin/app-macos/

test:
	go test -v ./... -count=1

conda:
	conda init
	conda activate go

requirements:
	go mod tidy