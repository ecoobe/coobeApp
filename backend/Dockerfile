# --- СТАДИЯ 1: сборка бинарника ---
FROM golang:1.20-alpine AS build

WORKDIR /src

# Настройка DNS и репозиториев Alpine
RUN echo "hosts: files dns" > /etc/nsswitch.conf \
	&& echo "http://dl-cdn.alpinelinux.org/alpine/v3.19/main" > /etc/apk/repositories \
	&& echo "http://dl-cdn.alpinelinux.org/alpine/v3.19/community" >> /etc/apk/repositories

# Устанавливаем git с повторными попытками
RUN set -eux; \
	for i in 1 2 3; do apk add --no-cache git && break || sleep 3; done

# Копируем зависимости
COPY go.mod go.sum ./

# Настройки Go (только IPv4)
ENV GODEBUG=netdns=go \
	GOPROXY=https://proxy.golang.org,direct

RUN go mod download

# Копируем весь код
COPY . .

# Собираем бинарь
RUN go build -o /app/server ./cmd/server

# --- СТАДИЯ 2: рантайм ---
FROM alpine:3.18

WORKDIR /app

# Настройка DNS + репо
RUN echo "hosts: files dns" > /etc/nsswitch.conf \
	&& echo "http://dl-cdn.alpinelinux.org/alpine/v3.18/main" > /etc/apk/repositories \
	&& echo "http://dl-cdn.alpinelinux.org/alpine/v3.18/community" >> /etc/apk/repositories

# Устанавливаем сертификаты с повторными попытками
RUN set -eux; \
	for i in 1 2 3; do apk add --no-cache ca-certificates && break || sleep 3; done \
	&& update-ca-certificates

COPY --from=build /app/server .

EXPOSE 8080

CMD ["./server"]