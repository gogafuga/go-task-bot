# Go Task Manager Bot

Telegram-бот для управления задачами, написанный на Go с использованием PostgreSQL и Docker Compose.

## Возможности

- `/start` — запуск бота
- `/add Название задачи` — добавление задачи
- `/list` — просмотр списка задач
- `/done ID` — отметка задачи как выполненной
- `/delete ID` — удаление задачи
- Хранение задач в PostgreSQL
- Разделение задач по пользователям

## Технологии

- Go
- PostgreSQL 17
- Docker Compose
- Telegram Bot API
- pgx/v5

## Запуск

Для работы необходимы Go, Docker и Docker Compose.

1. Клонируйте репозиторий.
2. Создайте файл `.env` с переменными `BOT_TOKEN`, `DATABASE_URL` и `POSTGRES_PASSWORD`.
3. Запустите PostgreSQL через Docker Compose.
4. Запустите Go-приложение.

## О проекте

Учебный проект для практики Go backend-разработки, работы с HTTP API, SQL и контейнеризацией.