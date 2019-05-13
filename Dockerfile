FROM golang:1.12 as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN CGO_ENABLED=0 go build -o surfboard-exporter -mod vendor .

FROM alpine
COPY --from=builder /build/surfboard-exporter /usr/bin
CMD ["surfboard-exporter"]
