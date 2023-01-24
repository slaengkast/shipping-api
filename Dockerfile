FROM golang:bullseye AS builder

RUN apt update && apt install git
RUN useradd -u 10001 shipping-api

WORKDIR $GOPATH/src/shipping-api
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /usr/local/bin/shipping-api ./cmd/...

FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/local/bin/shipping-api /shipping-api
USER shipping-api

ENTRYPOINT ["/shipping-api"]
