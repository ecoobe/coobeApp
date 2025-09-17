FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY backend/ ./backend/

RUN go build -o /server ./backend/cmd/server

EXPOSE 8080

CMD ["/server"]