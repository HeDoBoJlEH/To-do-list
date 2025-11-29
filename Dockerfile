FROM golang:1.24.3-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o ./bin/app ./app/cmd/main.go

EXPOSE 8080

CMD ["./bin/app"]