# Build container
FROM golang:1.13-alpine as builder

RUN apk add --no-cache git

WORKDIR /app/

# Install deps
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build
RUN go build -o ./out/dnschecker .

# Run container
FROM alpine:3

WORKDIR /app

COPY --from=builder /app/out/dnschecker .

EXPOSE 8080

CMD ["./dnschecker"]

