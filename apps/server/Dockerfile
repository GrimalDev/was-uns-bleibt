FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/wub-server .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates
COPY --from=builder /out/wub-server /usr/local/bin/wub-server

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/wub-server"]
