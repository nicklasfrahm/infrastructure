FROM golang:1.21.3 AS build
ARG VERSION
ARG TARGET

WORKDIR /app
COPY go.* ./
RUN go mod download

ADD Makefile .
ADD api/ api/
ADD cmd/ cmd/
ADD pkg/ pkg/
RUN TARGET=$TARGET VERSION=$VERSION BINARY=app make build

FROM gcr.io/distroless/base:nonroot AS run
WORKDIR /
COPY --from=build /app/app .
USER 65532:65532
ENTRYPOINT [ "/app" ]
LABEL org.opencontainers.image.source https://github.com/nicklasfrahm/infrastructure
