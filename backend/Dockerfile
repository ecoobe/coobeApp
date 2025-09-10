# --- СТАДИЯ 1: сборка бинарника ---
FROM golang:1.20-alpine AS build

# Рабочая директория внутри контейнера
WORKDIR /src

# Устанавливаем нужные пакеты (git нужен для go modules)
RUN apk add --no-cache git

# Переносим файлы зависимостей (go.mod и go.sum)
COPY go.mod go.sum ./

# Настройки чтобы Go использовал IPv4 и стабильный proxy
ENV GODEBUG=netdns=go \
	GOPROXY=https://proxy.golang.org,direct \
	GONOSUMDB=* \
	GOINSECURE=*

# Загружаем зависимости
RUN go mod download

# Переносим весь исходный код
COPY . .

# Собираем бинарник
RUN go build -o /app/server ./cmd/server

# --- СТАДИЯ 2: минимальный образ для запуска ---
FROM alpine:3.18

# Рабочая директория
WORKDIR /app

# Устанавливаем сертификаты для HTTPS (иначе Go клиент может не ходить по HTTPS)
RUN apk add --no-cache ca-certificates

# Копируем бинарник из стадии сборки
COPY --from=build /app/server .

# Указываем порт
EXPOSE 8080

# Запускаем сервер
CMD ["./server"]