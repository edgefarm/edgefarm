# Tools that are used for e2e tests

## e2epub

Publishes messages into edgefarm network via dapr

### Build
```bash
docker buildx build -f test/tools/e2epub/Dockerfile --platform linux/arm64,linux/amd64 --push -t ci4rail/dev-e2epub .
```

