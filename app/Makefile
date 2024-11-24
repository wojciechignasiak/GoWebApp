build:
	go build -o bin/app cmd/application/main.go

run: build
	./bin/app

test:
	go test -v ./... -count=1

conda:
	conda init
	conda activate go

requirements:
	go mod tidy