FROM golang:1.21-alpine

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код бэкенда
COPY backend/ ./backend/

# Собираем приложение
RUN go build -o /server ./backend/cmd/server

EXPOSE 8080

CMD ["/server"]