FROM golang:1.24-alpine
WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
COPY backend/ ./backend/

## Сборка с использованием vendor (если он есть)
RUN go build -o /server ./backend/cmd/server
EXPOSE 8080
CMD ["/server"]