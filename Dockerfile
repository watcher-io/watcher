FROM golang:latest as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o watcher ./main.go

FROM ubuntu:latest
WORKDIR /root/
COPY --from=builder /app/watcher  .
COPY --from=builder /app/.env  .
EXPOSE 8080
CMD ["./watcher"]