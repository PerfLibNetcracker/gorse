FROM golang:1.15-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN mkdir -p /app/bin/data
RUN go mod download
COPY . .
RUN go build -o main .

FROM golang:1.15-alpine
RUN mkdir database
COPY --from=builder /app/main .
COPY --from=builder /app/bin/config/gorse_docker.toml .
EXPOSE 8081
CMD ["./main","serve","-c","gorse_docker.toml"]
