# Сервис динамического сегментирования пользователей
Репозиторий содержит решение тестового задания на позицию стажёра Backend.

## Задача
Требуется реализовать сервис, хранящий пользователя и сегменты, в которых он состоит (создание, изменение, удаление сегментов, а также добавление и удаление пользователей в сегмент). 
Полное описание задания находится в [ASSIGNMENT.md](ASSIGNMENT.md)

### Стэк технологий
- Golang v1.20
- PostgreSQL v15
- Роутер                - [labstack/echo](https://github.com/labstack/echo) 
- Валидация             - [go-playground/validator](https://github.com/go-playground/validator)
- Драйвер БД            - [lib/pq](https://github.com/lib/pq)
- Парсер конфигурации   - [ilyakaznacheev/cleanenv](https://github.com/ilyakaznacheev/cleanenv)

## API
- `/create` - Создание нового сегмента
- `/delete` - Удаление нового сегмента
- `/experiments` - Добавление/удаление пользователя в сегмент
- `/list` - Получение списка сегментов пользователя
  
Более полное описание API с примерами запросов можно посмотреть в [соответствующем OpenAPI файле](api/openapi.yaml).

## Конфигурация
### Переменные окружения
- `AVITO_DATABASE_DSN` - Имя источника данных для подключения

## Запуск
Для запуска приложения необходимо инициализировать базу данных, таблицы и docker volume:
```bash
    # Скачать репозиторий
    git clone https://github.com/psxzz/backend-trainee-assignment-2023.git

    # Запустить dev-среду
    docker compose -f docker-compose.dev.yml up --detach

    # Выполнить скрипты создания таблиц
    docker exec -i <db-container> psql -U postgres -d experimental_segments < ./migrations/*.sql
    
    # Отключить dev-среду
    docker compose -f docker-compose.dev.yml down
```
После этого можно запускать само приложение:
```bash
    docker compose -f docker-compose.yml up --detach
```

## Ход работы
- Реализован сервис с базовым функционалом в соответствии с заданием
- 