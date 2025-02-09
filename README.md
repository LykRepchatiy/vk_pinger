# vk_pinger
Были реализованы 4 контейнера со следующим функцианалом:
  - BACKEND-сервис обеспечивает взаимодействие FRONTEND-сервиса с базой данных PostgreSQL и сохраняет список состояний и IP адрессов в базу данных
  - FRONTEND-сервис позволяет получить информацию для просмотра состояний контейнеров с помощью TS + React + Vite
  - PINGER-сервис узнает состояния контейнеров через DockerAPI и отправляет информаию в BACKEND-сервис, так же был реализован функцианал просмотра нескольких сетей
  - База данных PostgreSQL хранит в себе IP адресса, состояния, время проверки состояния, дату последнего успешного пинга
## Инструкция по запуску:
  Для запуска приложения требуется:
  - Склонировать репозиторий
  - Перейти в директроию с проектом в src
  - Выполнить команду находясь в src:
    ```bash
    docker-compose up
  - Открыть браузер и подключиться по адресу:
    ```
     http://localhost:3000
