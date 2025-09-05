FROM golang:1.20-alpine AS build

WORKDIR /src

# Копируем и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /src/cmd/server
RUN go build -o /app/server .

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /app/server /server

EXPOSE 8080
CMD ["/server"]