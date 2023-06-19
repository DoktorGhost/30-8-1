# Установка базового образа
FROM golang:latest

# Установка рабочей директории внутри контейнера
WORKDIR /app

# Копирование файлов Go modules
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода в контейнер
COPY . .

# Установка переменных окружения для подключения к PostgreSQL
ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=qwerty123
ENV DB_NAME=testdb

# Сборка Go приложения
RUN go build -o app

# Установка команды для запуска исполняемого файла
CMD ["./app"]