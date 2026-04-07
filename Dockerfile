# 1. Образ с Go
FROM golang:1.21-alpine

# 2. Создаем рабочую папку внутри контейнера
WORKDIR /app

# 3. Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# 4. Копируем весь остальной код
COPY . .

# 5. Собираем бинарный файл
RUN go build -o main ./cmd/main.go

# 6. Запускаем программу
CMD ["./main"]