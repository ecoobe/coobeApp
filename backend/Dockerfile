# ---- Стадия сборки ----
FROM golang:1.20-alpine AS build

WORKDIR /src

RUN apk add --no-cache git

# Копируем go.mod/go.sum и качаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем код backend
COPY backend/ ./backend/

# Собираем бинарник
WORKDIR /src/backend/cmd/server
RUN go build -o /app/server .

# ---- Финальный образ ----
FROM alpine:3.18

WORKDIR /app

# Бинарник
COPY --from=build /app/server .

# Фронт кладём внутрь образа бэка (чтобы отдавать статику тоже через Go)
COPY frontend/ ./frontend/

EXPOSE 8080

CMD ["./server"]