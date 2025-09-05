FROM golang:1.20-bullseye AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /src/cmd/server
RUN go build -o /app/server .

FROM debian:bullseye-slim

# Установка зависимостей и настройка сети
RUN apt-get update && \
	apt-get install -y ca-certificates && \
	rm -rf /var/lib/apt/lists/* && \
	echo "nameserver 8.8.8.8" > /etc/resolv.conf && \
	echo "nameserver 8.8.4.4" >> /etc/resolv.conf

COPY --from=build /app/server /server

EXPOSE 8080
CMD ["/server"]