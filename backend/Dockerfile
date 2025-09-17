FROM golang:1.24-alpine

WORKDIR /app

# Копируем файлы зависимостей и загружаем модули
COPY go.mod go.sum ./
RUN GOPROXY=https://proxy.golang.org,direct go mod download

COPY backend/ ./backend/

RUN go build -o /server ./backend/cmd/server

EXPOSE 8080

CMD ["/server"]
