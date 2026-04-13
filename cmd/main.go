package main

import (
	"context"
	//"fmt"
	"log"
	"log/slog" // Новый стандарт Go 1.21
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	//"todo-proj/internal/database"
	"todo-proj/internal/handlers"
	"todo-proj/internal/service"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Настройка логирования (JSON формат удобен для Docker)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 2. Инициализация конфига
	cfg := loadConfig()

	// 3. Инициализация БД  (теперь с retry)
	dbpool := setupDatabase(cfg.dbURL)
	defer dbpool.Close()

	// 4. Сборка слоев приложения
	taskSvc := service.NewTaskService(dbpool)
	h := &handlers.Handler{Service: taskSvc}
	router := handlers.NewRouter(h)

	// 5. Запуск сервера
	srv := &http.Server{
		Addr:    cfg.port,
		Handler: router,
	}

	go func() {
		slog.Info("Cервер запущен", "addr", cfg.port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("ошибка сервера", "err", err)
		}
	}()

	// 6. Graceful Shutdown
	waitForShutdown(srv)	
}

// --- Вспомогательные функции для чистоты main ---
type config struct {
	dbURL string
	port  string
}

func loadConfig() config {
	if err := godotenv.Load(); err != nil {
		slog.Info(".env не найден, используем системные переменные")
	}
	
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL не установлена")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ":8080"
	}
	return config{dbURL: dbURL, port: port}
}

func setupDatabase(connStr string) *pgxpool.Pool {
	var err error
	maxRetries := 5
	delay := 2 * time.Second

	for i := 1; i <= maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		pool, err := pgxpool.Connect(ctx, connStr)
		if err == nil {
			err := pool.Ping(ctx) // Проверяем реальную связь

			if err == nil {
				cancel() // Успех!
				slog.Info("успешное подключение к Postgres")
				return pool
			}
		}

		// Если мы здесь, значит была ошибка (Connect или Ping)
		cancel() // <--- Закрываем вручную при неудаче перед time.Sleep
		slog.Warn("База еще не готова", 
			"attempt", i, 
			"max_attempts", maxRetries, 
			"err", err)

		if i < maxRetries {
			time.Sleep(delay)
		}
	}

	log.Fatalf("не удалось подключиться к БД: %v", err)
	return nil
}

func waitForShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("завершение работы сервера...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("ошибка при остановке", "err", err)
	}
	slog.Info("сервер остановлен")
}
