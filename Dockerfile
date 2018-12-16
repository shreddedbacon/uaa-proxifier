FROM golang AS builder
RUN go version
COPY . /go/src/github.com/shreddedbacon/uaa-proxifier/
WORKDIR /go/src/github.com/shreddedbacon/uaa-proxifier/
RUN set -x && \
    go get -v .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o uaaproxy .

# actual container
FROM alpine:3.7
RUN apk --no-cache add ca-certificates openssl
WORKDIR /root/
# generate self signed for testing
RUN openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=localhost"  -keyout server.key  -out server.crt
# bring the actual executable from the builder
COPY --from=builder /go/src/github.com/shreddedbacon/uaa-proxifier .
ENV UAA_URL=https://192.168.50.6:8443
ENV PROXY_SSL_CERT=server.crt
ENV PROXY_SSL_KEY=server.key
ENV SKIP_INSECURE=true
ENV PORT=8080
ENTRYPOINT ["./uaaproxy"]
