FROM golang:latest as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o watcher ./main.go

FROM alpine:latest
WORKDIR /app/
COPY --from=builder /app/watcher  .
COPY --from=builder /app/.env.prod  .
EXPOSE 8080
CMD ["./watcher","prod"]