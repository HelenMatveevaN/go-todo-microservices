FROM golang:1.21-alpine

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем код
COPY . .

# Собираем (укажите путь к вашему main.go, если он в корне, то просто .)
RUN CGO_ENABLED=0 GOOS=linux go build -o notifier-app .

CMD ["./notifier-app"]