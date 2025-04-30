FROM golang:1.24-alpine

WORKDIR /app

COPY ./common /app/common
COPY ./auth-service /app/auth-service


COPY ./auth-service/go.mod ./auth-service/go.sum /app/
COPY go.work /app/go.work
RUN go work sync


RUN go mod tidy

RUN go build -o auth-service /app/auth-service/cmd/main.go

EXPOSE 8081

CMD ["./auth-service"]
