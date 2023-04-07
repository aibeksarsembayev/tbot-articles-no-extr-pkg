# build
FROM golang:1.19 AS build

WORKDIR /app

COPY .  .

RUN CGO_ENABLED=0 GOOS=linux go build -o sbertbot .

# deploy
FROM alpine:latest

WORKDIR /

COPY --from=build /app .

CMD ["./sbertbot"]