FROM golang:1.20 AS build
ARG VERSION
ARG TARGET

WORKDIR /app
COPY go.* ./
RUN go mod download

ADD Makefile .
ADD cmd/ cmd/
ADD pkg/ pkg/
RUN TARGET=$TARGET VERSION=$VERSION BINARY=controller make build

FROM gcr.io/distroless/base:nonroot AS run
WORKDIR /
COPY --from=build /app/controller .
USER 65532:65532
ENTRYPOINT [ "/controller" ]
