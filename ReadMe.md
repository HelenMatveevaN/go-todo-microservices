# To-Do Microservices System (Go + gRPC + Redis)

Простое и эффективное приложение для управления списком задач с использованием слоистой архитектуры (Clean Architecture).Полноценная экосистема для управления задачами, построенная на микросервисной архитектуре с фокусом на производительность и надежность.

## Стек технологий
**Language**: Go 1.21+
**Database**: PostgreSQL (Driver: pgx/v4)
**Caching**: Redis (Pattern: Cache Aside)
**Communication**: gRPC (Protobuf)
**HTTP Router**: Chi (Production-ready routing)
**Architecture**: Clean Architecture (Handlers -> Service -> Database)
**Logs**: slog (Structured JSON logging)
**Infrastructure**: Docker, Docker Compose
**Configuration**: Cleanenv (Support for .env and system env vars)

## Ключевые особенности
- **High Performance**: Кеширование списка задач в Redis снижает нагрузку на Postgres и ускоряет GET-запросы.
- **Smart Caching**: Автоматическая инвалидация (очистка) кеша при создании, обновлении или удалении задач.
- **Resilience**: Логика Retry при подключении к БД и Redis (приложение дождется их готовности).
- **Observability**: Внедрены Middleware (RequestID, Logger, Recoverer) для удобного трекинга запросов и защиты от паник.
- **Graceful Shutdown**: Безопасная остановка всех сервисов без потери данных и активных соединений.

## Межсервисное взаимодействие
Система состоит из двух сервисов. При создании новой задачи _ToDo Service_ выступает в роли gRPC-клиента и отправляет уведомление в фоновой горутине в _Notifier Service_, обеспечивая максимальный отклик API.

## Быстрый запуск
Всё окружение (App, Postgres, Redis, Notifier) поднимается одной командой:
   ```bash
   docker-compose up --build
   ```
Интерфейс приложения доступен по адресу: http://localhost:8080

## API документация
Все запросы проходят через систему Middleware. Основные эндпоинты:

*   **Метод	Путь	Описание**
*   **GET**	/tasks	Получить список задач (с поддержкой кеша Redis)
*   **POST**	/tasks	Создать задачу (инвалидирует кеш + gRPC уведомление)
*   **PATCH**	/tasks/{id}	Обновить статус задачи (инвалидирует кеш)
*   **DELETE**	/tasks/{id}	Удалить задачу (инвалидирует кеш)
*   **GET**	/health	Healthcheck сервиса

## Миграции
Для управления схемой БД используется инструмент **Goose**.
Миграции применяются автоматически при старте или вручную:
   ```bash
   goose -dir migrations postgres "user=... dbname=..." up
   ```

## Тестирование
Бизнес-логика (слой Service) покрыта Unit-тестами с использованием моков:
   ```bash
   go test ./internal/service/...
   ```

## 📄 Лицензия
Этот проект создан в целях исследования программного продукта и доступен для свободного использования.
