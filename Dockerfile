# build
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app/

ENV PORT="8080"
ENV APP_ENV="PROD"
ENV DATABASE_NAME="backend_template"
RUN go build -o main ./cmd/app/main.go

# run
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080

CMD [ "./main" ]