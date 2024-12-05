FROM golang:1.23.4-alpine as builder
ARG BUILD_VERSION
ARG BUILD_TIME

COPY . /src/
WORKDIR /src
RUN mkdir /output && \
    go build -o /output/ynab-exporter -ldflags "-X 'github.com/mcbobke/ynab-exporter/cmd/ynab-exporter/version.BuildTime=${BUILD_TIME}' -X 'github.com/mcbobke/ynab-exporter/cmd/ynab-exporter/version.BuildVersion=${BUILD_VERSION}'" ./cmd/ynab-exporter

FROM alpine:3.15.0
COPY --from=builder /output/ynab-exporter /ynab-exporter
ENTRYPOINT [ "/ynab-exporter" ]
