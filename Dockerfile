FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o feedbacklab ./cmd/feedbacklab
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o healthcheck ./cmd/healthcheck

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=builder /app/feedbacklab /app/feedbacklab
COPY --from=builder /app/healthcheck /app/healthcheck
COPY --from=builder /app/migrations ./migrations

USER nonroot:nonroot

EXPOSE 8080 8081

ENTRYPOINT ["/app/feedbacklab"]
