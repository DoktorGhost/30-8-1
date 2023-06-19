# Установка базового образа
FROM golang:latest

# Установка PostgreSQL
RUN apt-get update && apt-get install -y postgresql postgresql-contrib

# Установка Git (необходимо для загрузки зависимостей)
RUN apt-get install -y git

# Копирование исходного кода в контейнер
COPY . /app
WORKDIR /app

# Загрузка зависимостей
RUN go get -d -v ./...
RUN go install -v ./...

# Установка переменных окружения для подключения к PostgreSQL
ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=qwerty123
ENV DB_NAME=testdb

# Запуск контейнера с PostgreSQL и выполнение тестов
CMD service postgresql start && sleep 5 && go test -v ./...