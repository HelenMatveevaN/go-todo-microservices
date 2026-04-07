package database

func CreateTask(db *pgxpool.Pool, title string) error {
	query := `INSERT INTO tasks (title, is_done) VALUES ($1, false)`

	_, err := db.Exec(context.Background(), query, title)
	if err != nil {
		return fmt.Errorf("ошибка при создании задачи '%s': %w", title, err)
	}

	fmt.Printf("Задача '%s' успешно добавлена!\n", title)
	return nil
}

// достает все задачи из базу и возвращает их в виде слайса
func GetTasks(db *pgxpool.Pool) ([]Task, error) {
	//выполняем запрос
	rows, err := db.Query(context.Background(), "SELECT id, title, is_done FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("не удалось получить задачи: %w", err)
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var t Task
		//Scan копирует данные из колонок таблицы в поля структуры
		err := rows.Scan(&t.ID, &t.Title, &t.IsDone)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении строки: %w", err)
		}
		tasks = append(tasks, t) //добавляем задачу в список

	}
	return tasks, nil
}