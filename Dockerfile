FROM golang:1.24.1-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app/

COPY src ./src
WORKDIR /app/src

RUN go mod download

ENV CGO_ENABLED=1

RUN go build -o app

FROM alpine:3.19

# Install SQLite runtime library only
RUN apk add --no-cache sqlite-libs

WORKDIR /app/
COPY --from=builder /app/src/app .

EXPOSE 12345

CMD ["./app"]