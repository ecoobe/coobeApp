FROM golang:1.24

# Устанавливаем DNS-серверы Google
RUN echo "nameserver 8.8.8.8" > /etc/resolv.conf && \
	echo "nameserver 8.8.4.4" >> /etc/resolv.conf

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY backend/ ./backend/

RUN go build -o /server ./backend/cmd/server

EXPOSE 8080

CMD ["/server"]