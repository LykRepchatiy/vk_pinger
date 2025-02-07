# # Этап сборки (builder stage)
# FROM golang:1.22.2 AS builder
# WORKDIR /

# # Копируем файлы модулей и загружаем зависимости (кэширование зависимостей)
# COPY go.mod go.sum ./
# RUN go mod download

# # Копируем исходный код
# COPY . .

# # Собираем приложение, получаем бинарник с именем main (можете назвать его как угодно)
# RUN go build -o backend .

# # Этап запуска (runtime stage)
# FROM alpine:latest
# WORKDIR /

# # Копируем бинарник из этапа сборки
# COPY --from=builder /backend .

# # Опционально, если ваше приложение использует TLS или другие сертификаты,
# # установите необходимые пакеты, например, ca-certificates:
# # RUN apk add --no-cache ca-certificates

# # Указываем команду для запуска приложения
# CMD ["./backend"]



# Базовый образ с Go
FROM golang:1.22.2

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта в контейнер
COPY . .

# Компилируем приложение
RUN go build -o backend .

# Указываем порт для прослушивания
EXPOSE 8080

# Указываем команду для запуска приложения
CMD ["./backend"]