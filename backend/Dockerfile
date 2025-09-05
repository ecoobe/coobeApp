FROM golang:1.20-alpine AS build

# Устанавливаем инструменты для работы с DNS
RUN apk add --no-cache ca-certificates

WORKDIR /src

COPY go.mod go.sum ./

# Используем прямой доступ к репозиториям (без proxy)
ENV GOPROXY=direct
ENV GOSUMDB=off

RUN go mod download

COPY . .

WORKDIR /src/cmd/server
RUN go build -o /app/server .

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /app/server /server

EXPOSE 8080
CMD ["/server"]