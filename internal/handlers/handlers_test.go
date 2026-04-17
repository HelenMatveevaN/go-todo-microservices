package handlers

import (
	"context"
	//"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo-proj/internal/models"
	"todo-proj/internal/service" // Добавили service

	"github.com/go-chi/chi/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

// 1. Мок-сервис
type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) List(ctx context.Context) ([]models.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskService) GetByID(ctx context.Context, id int) (models.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockTaskService) Create(ctx context.Context, title string) (models.Task, error) {
	args := m.Called(ctx, title)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockTaskService) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTaskService) UpdateStatus(ctx context.Context, id int, isDone bool) error {
	args := m.Called(ctx, id, isDone)
	return args.Error(0)
}

// ТЕСТ 1: Успешный список
func TestGetTasksHandler(t *testing.T) {
	mockSvc := new(MockTaskService)
	h := &Handler{Service: mockSvc}

	mockData := []models.Task{{ID: 1, Title: "Тест"}}
	mockSvc.On("List", mock.Anything).Return(mockData, nil)

	req := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	h.GetTasksHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

// ТЕСТ 2: Ошибка 404
func TestGetTaskByID_NotFound(t *testing.T) {
	mockSvc := new(MockTaskService)
	h := &Handler{Service: mockSvc}

	mockSvc.On("GetByID", mock.Anything, 999).Return(models.Task{}, service.ErrTaskNotFound)

	req := httptest.NewRequest("GET", "/tasks/999", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	
	w := httptest.NewRecorder()

	h.GetTaskByIDHandler(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}