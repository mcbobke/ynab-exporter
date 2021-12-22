FROM golang:1.17-alpine as builder
COPY . /src/
WORKDIR /src
RUN mkdir /output && \
    go build -o /output/ynab-exporter ./cmd/ynab-exporter

FROM alpine:3.15.0
COPY --from=builder /output/ynab-exporter /ynab-exporter
ENTRYPOINT [ "/ynab-exporter" ]
