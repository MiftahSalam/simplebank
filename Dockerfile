FROM golang:1.20.0-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o simplebank main.go
RUN apk add curl
RUN apk add --no-cache bash
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine
WORKDIR /app
COPY --from=builder /app/simplebank .
COPY --from=builder /app/migrate .
COPY db/migrations ./migrations
COPY app.env .
COPY wait-for.sh .
COPY start.sh .

EXPOSE 8080
CMD [ "/app/simplebank" ]
ENTRYPOINT [ "/app/start.sh" ]