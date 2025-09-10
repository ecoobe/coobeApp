# --- СТАДИЯ 1: сборка бинарника ---
FROM golang:1.20-alpine AS build

WORKDIR /src

# Настройка: отключаем IPv6, используем IPv4
RUN echo "hosts: files dns" > /etc/nsswitch.conf \
	&& echo "http://dl-cdn.alpinelinux.org/alpine/v3.19/main" > /etc/apk/repositories \
	&& echo "http://dl-cdn.alpinelinux.org/alpine/v3.19/community" >> /etc/apk/repositories

# Устанавливаем git
RUN apk add --no-cache git

# Копируем зависимости
COPY go.mod go.sum ./

# Настройки для go mod
ENV GODEBUG=netdns=go \
	GOPROXY=https://proxy.golang.org,direct

RUN go mod download

# Копируем код
COPY . .

# Собираем бинарь
RUN go build -o /app/server ./cmd/server

# --- СТАДИЯ 2: рантайм ---
FROM alpine:3.18

WORKDIR /app

# Настройка: IPv4-only, обновление репозиториев
RUN echo "hosts: files dns" > /etc/nsswitch.conf \
	&& echo "http://dl-cdn.alpinelinux.org/alpine/v3.18/main" > /etc/apk/repositories \
	&& echo "http://dl-cdn.alpinelinux.org/alpine/v3.18/community" >> /etc/apk/repositories \
	&& apk add --no-cache ca-certificates \
	&& update-ca-certificates

COPY --from=build /app/server .

EXPOSE 8080

CMD ["./server"]