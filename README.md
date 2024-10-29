# sentry-exporter

## Build Docker image
```bash
docker build -t jushcherbak/sentry-exporter:$(cat VERSION.txt) --build-arg BUILD_VERSION=$(cat VERSION.txt) -f docker/Dockerfile .
```
pulling image:
```bash
jushcherbak/sentry-exporter:0.0.2
```

## Usage

### secure `/metrics` with Bearer token

```bash
export API_TOKEN=secret
go run main.go
```

```bash
curl http://127.0.0.1:9101/metrics -H "Authorization: Bearer secret"
```

### filter exported projects by filtes

```bash
export PROJECTS_EXCLUDE=".+(\_|\-)dev$"
go run main.go
```

If `PROJECTS_INCLUDE` and `PROJECTS_EXCLUDE` are set the applied in chain: include by regex, exclude by regex