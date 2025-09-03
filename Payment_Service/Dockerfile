# ---------- build stage ----------
FROM golang:1.24.1 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/payment-api ./main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/outboxrelay ./cmd/outboxrelay/main.go

# ---------- runtime stage ----------
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/payment-api /app/payment-api
COPY --from=build /out/outboxrelay /app/outboxrelay
COPY ./env ./env
USER 65532:65532
ENTRYPOINT ["/app/payment-api"]
