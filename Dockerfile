# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.19 as build-stage

WORKDIR /app

COPY . ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/web ./cmd/web

# RUN the tests in the container
# FROM build-stage as run-test-stage

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 as build-release-stage

WORKDIR /app

# copy ui template artifacts from build-stage
COPY --from=build-stage /app/ui ./ui

# copy executable artifacts from build-stage
COPY --from=build-stage /app/bin/web ./bin/web

USER nonroot:nonroot

EXPOSE 4000

CMD ["./bin/web"]

