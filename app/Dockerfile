FROM golang:1.23

WORKDIR /golang_api

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o bin/app cmd/application/main.go

CMD ["./bin/app"]
