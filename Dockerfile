FROM golang:alpine as builder
RUN mkdir /build
COPY . /build/
WORKDIR /build
RUN CGO_ENABLED=0 go test ./app/... &&\
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o joytotwi ./app/

FROM scratch
COPY --from=builder /build/joytotwi /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
ENTRYPOINT ["./joytotwi"]