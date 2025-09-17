FROM golang:1.21-alpine

# Устанавливаем необходимые пакеты и настраиваем DNS
RUN apk update && apk add --no-cache git curl && \
	echo "nameserver 8.8.8.8" > /etc/resolv.conf && \
	echo "nameserver 8.8.4.4" >> /etc/resolv.conf

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Пробуем разные подходы к загрузке модулей
RUN GOPROXY=direct GOSUMDB=off go mod download || \
	GOPROXY=https://proxy.golang.org GOSUMDB=sum.golang.org go mod download || \
	GOPROXY=https://goproxy.io GOSUMDB=off go mod download

# Копируем и собираем приложение
COPY backend/ ./backend/
RUN go build -o /server ./backend/cmd/server

EXPOSE 8080

CMD ["/server"]