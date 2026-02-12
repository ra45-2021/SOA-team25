FROM golang:1.24 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o service ./cmd/api


FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=build /app/service .

EXPOSE 8080
CMD ["./service"]
