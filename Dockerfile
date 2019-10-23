FROM golang:alpine as tests
ARG SKIP_TEST
RUN mkdir /build
COPY . /build/
WORKDIR /build
RUN if [ -z "${SKIP_TEST}" ] ; then \
        CGO_ENABLED=0 go test ./app/...; \
    else \
        echo "skipping tests"; \
    fi

FROM tests as builder
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o joytotwi ./app/

FROM scratch
COPY --from=builder /build/joytotwi /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
ENTRYPOINT ["./joytotwi"]