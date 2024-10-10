FROM registry.access.redhat.com/ubi9/go-toolset:latest AS build

WORKDIR /opt/app-root/src
COPY . .

ENV GOOS=linux
ENV GOARCH=amd64

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w" -o /opt/app-root/src/price-service cmd/main.go

FROM redhat/ubi9-micro:latest

COPY --from=build /opt/app-root/src/price-service /usr/bin/price-service

CMD ["/usr/bin/price-service"]
