FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./api

FROM gcr.io/distroless/base-debian12

WORKDIR /running

COPY --from=builder /app/api ./api
COPY --from=builder /app/seed ./seed

EXPOSE 8080
CMD ["./api"]


