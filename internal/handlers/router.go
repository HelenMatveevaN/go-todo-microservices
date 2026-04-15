package handlers

import (
    "net/http"
    "time"

    "github.com/go-chi/chi/v4"
    "github.com/go-chi/chi/v4/middleware"
)

// NewRouter настраивает маршруты и промежуточное ПО (middleware)
func NewRouter(h *Handler) *chi.Mux {
    r := chi.NewRouter()

    // 1. Стандартные Middleware
    r.Use(middleware.RequestID)    // Добавляет ID к каждому запросу для трекинга
    r.Use(middleware.RealIP)       // Определяет настоящий IP пользователя
    r.Use(middleware.Logger)       // Логирует запросы
    r.Use(middleware.Recoverer)    // Защита от паник

    // Устанавливаем таймаут на обработку запроса
    r.Use(middleware.Timeout(60 * time.Second))

    // 2. Маршруты API
    r.Get("/health", HealthCheck) // Простая проверка доступности

    r.Route("/tasks", func(r chi.Router) {      
        r.Get("/", h.GetTasksHandler)           // GET /tasks
        r.Get("/{id}", h.GetTaskByIDHandler)    // GET /tasks/123
        r.Post("/", h.CreateTaskHandler)        // POST /tasks
        r.Delete("/{id}", h.DeleteTaskHandler)  // DELETE /tasks/123
        r.Patch("/{id}", h.UpdateTaskHandler)   // PATCH /tasks/123
    })

    // 3. Раздача статики (Frontend)
    fs := http.FileServer(http.Dir("./static"))
    r.Handle("/*", http.StripPrefix("/", fs))
    
    return r
}