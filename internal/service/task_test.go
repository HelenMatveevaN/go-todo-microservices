package service

import (
	"testing"
)

func TestValidateTask(t *testing.T) {
	// Описываем тестовые случаи (Table-driven tests)
	tests := []struct {
		name    string
		title   string
		wantErr error
	}{
		{
			name:    "Валидный заголовок",
			title:   "Купить молоко",
			wantErr: nil,
		},
		{
			name:    "Пустой заголовок",
			title:   "",
			wantErr: ErrTitleTooEmpty,
		},
		{
			name:    "Слишком длинный заголовок",
			title:   "Это очень очень очень очень очень очень длинный заголовок больше ста символов для проверки лимита",
			wantErr: ErrTitleTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Вызываем твою логику валидации
			err := ValidateTask(tt.title) 

			if err != tt.wantErr {
				t.Errorf("ValidateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}