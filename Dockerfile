FROM golang:1.12 as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go build -o surfboard_exporter .

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/surfboard_exporter /app/
WORKDIR /app
CMD ["./surfboard_exporter"]
