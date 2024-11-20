FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o analytics-platform cmd/app/main.go



FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/analytics-platform .

EXPOSE 8080

CMD ["./analytics-platform"]



# For db - pg_dump -U postgres -h localhost -p 5432 -d truck-analytics > data_dump.sql

#FROM postgres:14

#COPY data_dump.sql /docker-entrypoint-initdb.d/