################################################### STAGE 1 -- BUILD
FROM golang:1.22.3-alpine3.20 AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./main.go

################################################### STAGE 2 -- RUN
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/main ./

EXPOSE 8080

ENTRYPOINT ["./main"]