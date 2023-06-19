package storage

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

/*
+Создавать новые задачи,
+Получать список всех задач,
+Получать список задач по автору,
+Получать список задач по метке,
+Обновлять задачу по id,
Удалять задачу по id.
*/

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// Задача.
type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// Tasks возвращает список задач из БД.
func (s *Storage) Tasks(taskID, authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2)
		ORDER BY id;
	`,
		taskID,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// NewTask создаёт новую задачу и возвращает её id.
func (s *Storage) NewTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (author_id, title, content)
		VALUES ($1, $2, $3) RETURNING id;
		`,
		t.AuthorID,
		t.Title,
		t.Content,
	).Scan(&id)
	return id, err
}

// NewLabel создаёт новую запись в таблице labels (нееобходимо для тестирования)
func (s *Storage) NewLabel(name string) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO labels (name)
		VALUES ($1) RETURNING id;
		`,
		name,
	).Scan(&id)
	return id, err
}

// AddLabelToTask создаёт новую запись в таблице tasks_labels (нееобходимо для тестирования)
func (s *Storage) NewTasksLabels(taskID, labelID int) error {
	_, err := s.db.Exec(context.Background(), `
		INSERT INTO tasks_labels (task_id, label_id)
		VALUES ($1, $2);
		`,
		taskID,
		labelID,
	)
	if err != nil {
		return err
	}
	return nil
}

// TasksAuthorName возвращает список задач  определенного автора, принмает строку - имя автора
func (s *Storage) TasksAuthorName(name string) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE author_id IN (
			SELECT id 
			FROM users
			WHERE name = $1
		)
		ORDER BY id;
	`,
		name,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// роверить rows.Err()
	return tasks, rows.Err()
}

// TasksLabels возвращает список задач по метке
func (s *Storage) TasksLabels(name string) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE id IN (
			SELECT task_id
			FROM tasks_labels
			WHERE label_id IN (
				SELECT id 
				FROM labels
				WHERE name = $1
			)
		)
		ORDER BY id;
	`,
		name,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// роверить rows.Err()
	return tasks, rows.Err()
}

// correct изменеение данных по номеру ID задачи
func (s *Storage) TasksCorrect(t Task) ([]Task, error) {
	_, err := s.db.Exec(context.Background(), `
		UPDATE tasks
		SET closed = $2,
			assigned_id = $3,
			title = $4,
			content = $5
		WHERE id = $1;
	`,
		t.ID,
		t.Closed,
		t.AssignedID,
		t.Title,
		t.Content,
	)
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRow(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE id = $1;
	`,
		t.ID,
	)

	var task Task
	err = row.Scan(
		&task.ID,
		&task.Opened,
		&task.Closed,
		&task.AuthorID,
		&task.AssignedID,
		&task.Title,
		&task.Content,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Запись не найдена
		}
		return nil, err
	}

	return []Task{task}, nil
}

// DeleteTasks удаляет запись в таблице tasks
func (s *Storage) DeleteTasks(taskID int) error {
	_, err := s.db.Exec(context.Background(), `
		DELETE FROM tasks
		WHERE id = $1;
		`,
		taskID,
	)
	if err != nil {
		return err
	}
	return nil
}

// получение задачипо ID
func (s *Storage) GetTasks(taskID int) ([]Task, error) {
	row := s.db.QueryRow(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
		id = $1;
	`,
		taskID,
	)

	var task Task
	err := row.Scan(
		&task.ID,
		&task.Opened,
		&task.Closed,
		&task.AuthorID,
		&task.AssignedID,
		&task.Title,
		&task.Content,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Запись не найдена
		}
		return nil, err
	}

	return []Task{task}, nil
}
