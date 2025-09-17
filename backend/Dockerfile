FROM golang:1.24-alpine
WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
COPY backend/ ./backend/
## Копируем vendor, если он есть
COPY backend/vendor ./vendor
## Сборка с использованием vendor (если он есть)
RUN go build -mod=vendor -o /server ./backend/cmd/server
EXPOSE 8080
CMD ["/server"]