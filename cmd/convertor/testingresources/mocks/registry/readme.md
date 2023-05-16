# Local Registry
This package provides a local implementation of a registry complete with 3 sample images of different types. In particular this package is meant to provide a local registry for testing purposes. Note that the local registry is not particularly optimized or a good model for how to implement a local registry. In particular we are using abstractions from https://pkg.go.dev/github.com/containers/image/v5 for the purpose of maintaining compatibility with skopeo image downloads as a quick, easy and reproducible way of adding and downloading images.

## Present Images

### Docker
```bash
skopeo sync --src docker --dest dir mcr.microsoft.com/mcr/hello-world:latest ./registry
```

### List
```bash
skopeo sync --all --src docker --dest dir docker.io/library/redis:alpine3.18 ./registry
```