FROM golang:1.19.10-bullseye

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем исходный код в контейнер
COPY . .

# Компилируем приложение
RUN go build -o app

# Запускаем контейнер и оставляем его активным
CMD tail -f /dev/null