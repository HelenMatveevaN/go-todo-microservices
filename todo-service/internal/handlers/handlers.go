package handlers

import (
	"encoding/json"
	"errors"
	"log/slog" // Добавляем slog
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v4"

	"todo-proj/internal/service"
)

type Handler struct {
	Service service.TaskService //зависим от интерфейса
}

type Response struct {
	Data  interface{} `json:"data,omitempty"`  // Любые данные (объект или список)
	Error string      `json:"error,omitempty"` // Текст ошибки
}

// Хелперы для унификации ответов
func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{Data: data})
}

func sendError(w http.ResponseWriter, status int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Логируем ошибку вместе с HTTP статусом и текстом
	slog.Error("запрос завершился ошибкой", 
		"status", status, 
		"msg", message, 
		"err", err,
	)

	json.NewEncoder(w).Encode(Response{Error: message})
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, http.StatusOK, "API To-Do приложение работает!")
}

func (h *Handler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Service.List(r.Context())
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Ошибка БД", err)
		return
	}
	sendJSON(w, http.StatusOK, tasks)
}

func (h *Handler) GetTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id") // Вытаскиваем ID из ссылки
	id, err := strconv.Atoi(idStr) // Конвертируем "1" в число 1
	if err != nil {
		sendError(w, http.StatusBadRequest, "Некорректный ID", err)
		return
	}
	
	task, err := h.Service.GetByID(r.Context(), id)
	if err != nil {
		sendError(w, http.StatusNotFound, "Задача не найдена", err)
		return
	}
	sendJSON(w, http.StatusOK, task)
}

func (h *Handler) CreateTaskHandler (w http.ResponseWriter, r *http.Request) {
	var req struct { 
		Title string `json:"title"` 
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Некорректный JSON", err)
		return
	}

	task, err := h.Service.Create(r.Context(), req.Title)
	if err != nil {
		status := http.StatusInternalServerError
		// Проверяем конкретные ошибки сервиса
		if errors.Is(err, service.ErrTitleTooEmpty) || errors.Is(err, service.ErrTaskInvalidTitle) || errors.Is(err, service.ErrTitleTooLong) {
			status = http.StatusBadRequest
		}
		sendError(w, status, err.Error(), err)
		return
	}

	// Если ошибок нет, возвращаем созданную задачу
	slog.Info("создана новая задача через API", "id", task.ID)
	sendJSON(w, http.StatusCreated, task)
}

func (h *Handler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Неверный ID", err)
		return
	}

	var input struct {
		IsDone bool `json:"is_done"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Плохой JSON", err)
		return
	}

	err = h.Service.UpdateStatus(r.Context(), id, input.IsDone)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			sendError(w, http.StatusNotFound, err.Error(), err)
		} else {
			sendError(w, http.StatusInternalServerError, "Не удалось обновить статус", err)
		}
		return
	}
	sendJSON(w, http.StatusOK, map[string]string{"message": "Статус обновлен"})
}

func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id") // достаем {id} из URL
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Неверный ID", err)
		return
	}

	if err := h.Service.Delete(r.Context(), id); err != nil {
		sendError(w, http.StatusInternalServerError, "Ошибка удаления", err)
		return
	}
	sendJSON(w, http.StatusOK, map[string]string{"message": "Удалено"})
}

