FROM golang:1.24.1 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/product-catalog-service ./main.go

FROM  gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/product-catalog-service /app/product-catalog-service
COPY ./env ./env
ENTRYPOINT [ "/app/product-catalog-service" ]