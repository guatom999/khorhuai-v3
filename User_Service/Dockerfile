FROM golang:1.24.1 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/user-api

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/user-api /app/user-api
COPY ./env ./env
ENTRYPOINT [ "/app/user-api" ]