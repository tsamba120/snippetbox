# syntax=docker/dockerfile:1

FROM golang:1.19

WORKDIR /app

COPY . ./

RUN ls -la

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/web ./cmd/web

EXPOSE 4000

CMD ["./bin/web"]

