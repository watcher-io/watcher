FROM golang:latest as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CORS_ORIGIN ""
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o main ./main.go

FROM ubuntu:latest
WORKDIR /root/
COPY --from=builder /app/main  .
COPY --from=builder /app/.env  .
COPY --from=builder /app/cert/  ./cert
EXPOSE 8080
CMD ["./main"]