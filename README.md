# Go Microservices Ecosystem 🚀

Комплексная распределенная система на Go, реализующая событийную архитектуру, кеширование и межсервисное взаимодействие.

## 🏗 Архитектура системы
Система разделена на три функциональных микросервиса:
1.  **ToDo Service (Core)**: REST API для управления задачами. Работает с PostgreSQL и реализует кеширование через Redis. Выступает издателем (Publisher) событий в RabbitMQ.
2.  **Notifier Service (gRPC)**: Сервис уведомлений. Принимает прямые вызовы от Core-сервиса по протоколу gRPC.
3.  **Statistics Service (Event Consumer)**: Асинхронный воркер. Слушает очередь RabbitMQ и обрабатывает события создания задач в реальном времени.

## 🛠 Технологический стек
*   **Backend**: Go 1.21+ (Clean Architecture)
*   **Communication**: gRPC (Protobuf), RabbitMQ (AMQP 0.9.1)
*   **Databases**: PostgreSQL (pgx/v4), Redis (Cache Aside pattern)
*   **Observability**: Structured Logging (slog), Middlewares (RequestID, Recovery, Logger)
*   **DevOps**: Docker, Docker Compose (контейнеризация всей инфраструктуры)

## 🚀 Быстрый запуск
Вся система поднимается одной командой из корневой директории:

```bash
docker-compose up --build
```

## API Gateway: 
http://localhost:8080

## RabbitMQ Management: 
http://localhost:15672 (guest/guest)

## 🛡 Надежность (Resilience)

*   **Retry Logic**: Все сервисы имеют встроенную логику ожидания готовности инфраструктуры (DB, RabbitMQ) при холодном старте.
*   **Graceful Shutdown**: Сервисы корректно завершают работу, закрывая соединения с брокером и базами данных.
*   **Context Control**: Управление таймаутами gRPC и HTTP запросов.
  
## 📄 Разработка
Каждый сервис расположен в своей директории:
/todo-service — основная бизнес-логика и миграции.
/notifier-service — реализация gRPC сервера.
/statistics-service — обработчик событий.
