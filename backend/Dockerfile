# backend/Dockerfile
FROM golang:1.20-alpine AS build
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /src/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]