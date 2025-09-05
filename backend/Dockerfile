FROM golang:1.20-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./

# Пробуем использовать альтернативный DNS
RUN echo "nameserver 8.8.8.8" > /etc/resolv.conf && \
	echo "nameserver 8.8.4.4" >> /etc/resolv.conf && \
	go mod download

COPY . .

WORKDIR /src/cmd/server
RUN go build -o /app/server .

FROM scratch
COPY --from=build /app/server /server

EXPOSE 8080
CMD ["/server"]