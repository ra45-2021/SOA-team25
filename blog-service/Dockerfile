FROM golang:1.22 AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o service main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /
EXPOSE 8080
COPY --from=build /app/service /service
ENTRYPOINT ["/service"]
