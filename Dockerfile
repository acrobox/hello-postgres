FROM golang:alpine AS builder
ENV CGO_ENABLED=0
COPY . /build/
WORKDIR /build
RUN go build -a -installsuffix docker -ldflags='-w -s' -o /build/bin/hello-postgres /build

FROM ghcr.io/acrobox/docker/minimal:latest
EXPOSE 8080
COPY --from=builder /build/bin/hello-postgres /usr/local/bin/hello-postgres
USER user
CMD ["/usr/local/bin/hello-postgres"]

LABEL org.opencontainers.image.source https://github.com/acrobox/hello-postgres
