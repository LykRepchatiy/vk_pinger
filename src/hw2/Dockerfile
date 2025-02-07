# Используем базовый образ Alpine Linux
FROM alpine:latest

# Устанавливаем метку для описания образа
LABEL description="A container that prints 'Hello World' in an infinite loop"

EXPOSE 9090

# Команда, которая будет выполняться при запуске контейнера
CMD while true; do echo "Hello World"; sleep 1; done