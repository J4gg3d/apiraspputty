FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=arm GOARM=7 go build -o restapi ./server/main.go

FROM debian:rc-buggy-20240812

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /app

COPY --from=builder /app/restapi /app/restapi

# Expose the application port
EXPOSE 8080

ENTRYPOINT ["/app/restapi"]