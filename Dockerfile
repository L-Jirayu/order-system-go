FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app .

# Final Stage
FROM gcr.io/distroless/base

WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 8080
CMD ["./app"]
