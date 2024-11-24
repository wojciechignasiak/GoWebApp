FROM golang:1.23

WORKDIR /golang_api

# Kopiujemy pliki go.mod i go.sum
COPY go.mod go.sum ./

# Pobieramy zależności
RUN go mod download && go mod verify

# Kopiujemy całą resztę projektu
COPY . .

# Budujemy aplikację
RUN go build -o bin/app cmd/application/main.go

# Uruchamiamy aplikację
CMD ["./bin/app"]
