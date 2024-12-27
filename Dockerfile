FROM golang:1.23-alpine

WORKDIR /app

# Установка git
RUN apk add --no-cache git

# Копирование go.mod и go.sum (если есть)
COPY go.mod ./
COPY go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN go build -o bin/backend cmd/main.go

EXPOSE 8181

CMD ["./bin/backend"] 