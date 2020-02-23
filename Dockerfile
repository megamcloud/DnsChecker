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

CMD ["./out/dnschecker"]
