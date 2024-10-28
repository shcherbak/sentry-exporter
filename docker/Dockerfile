FROM golang:1.21.11-bullseye AS builder
ARG BUILD_VERSION='0.0.1'
ENV LISTEN_PORT=9101
WORKDIR /opt/sentry-exporter/
COPY go.mod ./
COPY go.sum .
COPY src ./src
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.buildVersion=${BUILD_VERSION}" -a -installsuffix cgo -o sentry-exporter src/main.go

FROM debian:bullseye
WORKDIR /opt/sentry-exporter/
RUN apt-get update \
  && apt-get install -y --no-install-recommends \
    curl \
  && rm -rf /var/cache/apt/archives/*
COPY --from=builder /opt/sentry-exporter/sentry-exporter .

EXPOSE ${LISTEN_PORT}
CMD ["/opt/sentry-exporter/sentry-exporter"]
