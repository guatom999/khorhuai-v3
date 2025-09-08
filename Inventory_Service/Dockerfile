FROM golang:1.24.1 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/inventory-api ./main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/worker ./corn/worker/main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/inventory-api /app/inventory-api
COPY --from=build /out/worker        /app/worker 
COPY ./env ./env