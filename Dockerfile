FROM golang AS builder
RUN mkdir /build
COPY . /build/
WORKDIR /build
RUN go get all
RUN CGO_ENABLED=0 GOOS=linux go build -o nilan2mqtt cmd/nilan2mqtt/main.go

FROM alpine:latest AS runtime
FROM scratch AS final
COPY --from=builder /build/nilan2mqtt .
ENTRYPOINT [ "./nilan2mqtt" ]
