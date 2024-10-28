# sentry-exporter

## Build Docker image
```bash
docker build -t jushcherbak/sentry-exporter:$(cat VERSION.txt) --build-arg BUILD_VERSION=$(cat VERSION.txt) -f docker/Dockerfile .
```
pulling image:
```bash
jushcherbak/sentry-exporter:0.0.1
```