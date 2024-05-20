FROM golang:1.22 as builder
WORKDIR /build

COPY . .
RUN rm config.yaml; exit 0

RUN make docker-build

FROM alpine:latest

COPY --from=builder build/ .
RUN addgroup --gid 1001 -S consumer-api && \
    adduser -G consumer-api --shell /bin/false --disabled-password -H --uid 1001 consumer-api
USER consumer-api
EXPOSE 6060
ENTRYPOINT [ "./consumer-api" ]