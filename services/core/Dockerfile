FROM golang:1.26.4-alpine AS builder

WORKDIR /app

RUN apk --no-cache add bash git make gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd/skilltracker/main.go

FROM alpine:latest

COPY --from=builder ./app/app ./

EXPOSE 8080
CMD ["./app"]