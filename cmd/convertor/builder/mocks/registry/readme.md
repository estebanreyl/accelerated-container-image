# Local Registry

## Images

### Docker

```bash
skopeo sync --src docker --dest dir mcr.microsoft.com/mcr/hello-world:latest ./registry
```

### List

```bash
skopeo sync --all --src docker --dest dir docker.io/library/redis:alpine3.18 ./registry
```