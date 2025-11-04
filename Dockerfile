FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o innotech ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /app/innotech /app/innotech
COPY --from=builder /app/migrations ./migrations

USER nonroot:nonroot

EXPOSE 8080
ENTRYPOINT ["/app/innotech"]
