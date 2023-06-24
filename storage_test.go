package storage_test

import (
	"storage"
	"testing"
)

// Создайте экземпляр хранилища
func connectDb(t *testing.T) *storage.Storage {
	s, err := storage.New("postgres://test_user:qwerty123@localhost:5432/testdb")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	return s
}

func TestTasks(t *testing.T) {
	s := connectDb(t)

	// Вызовите функцию Tasks с тестовыми аргументами
	taskID := 0
	authorID := 0
	tasks, err := s.Tasks(taskID, authorID)
	if err != nil {
		t.Fatalf("Failed to get tasks: %v", err)
	}
	t.Log(tasks)

	taskID = 1
	authorID = 0
	tasks, err = s.Tasks(taskID, authorID)
	if err != nil {
		t.Fatalf("Failed to get tasks: %v", err)
	}
	t.Log(tasks)

}

func TestNewTask(t *testing.T) {
	s := connectDb(t)

	// Создайте тестовую задачу
	task := storage.Task{
		AuthorID: 3,
		Title:    "Test Task",
		Content:  "Test Content",
	}

	// Вызовите функцию NewTask
	taskID, err := s.NewTask(task)
	if err != nil {
		t.Fatalf("Failed to create new task: %v", err)
	}

	t.Log("Creat task with id:", taskID)

}

func TestTasksAuthorName(t *testing.T) {
	s := connectDb(t)

	// Вызовите функцию TasksAuthorName с тестовыми аргументами
	name := "geralt"

	tasks, err := s.TasksAuthorName(name)
	if err != nil {
		t.Fatalf("Failed to get tasks: %v", err)
	}
	t.Log(tasks)

	name = "ollya"
	tasks, err = s.TasksAuthorName(name)
	if err != nil {
		t.Fatalf("Failed to get tasks: %v", err)
	}
	t.Log(tasks)

	name = "ollya2"
	tasks, err = s.TasksAuthorName(name)
	if err != nil {
		t.Fatalf("Failed to get tasks: %v", err)
	}
	t.Log(tasks)

}

func TestTasksLabels(t *testing.T) {
	s := connectDb(t)

	name := "LACOSTA"

	// Добавьте метку в базу данных
	labelID, err := s.NewLabel(name)
	if err != nil {
		t.Fatalf("Failed to add new label: %v", err)
	}
	t.Log("Add:", name)

	// Создайте новую задачу
	task := storage.Task{
		AuthorID: 1,
		Title:    "New Task with label",
		Content:  "Test Content with label",
	}

	// Добавьте задачу в базу данных
	taskID, err := s.NewTask(task)
	if err != nil {
		t.Fatalf("Failed to add new task: %v", err)
	}
	t.Log("Add:", task.Title)

	// Свяжите метку с задачей
	err = s.NewTasksLabels(taskID, labelID)
	if err != nil {
		t.Fatalf("Failed to add label to task: %v", err)
	}
	t.Log("Add connection:", taskID, "and", labelID)

	labels, err := s.TasksLabels(name)
	if err != nil {
		t.Fatalf("Failed to get labels: %v", err)
	}
	t.Log(labels)

}

func TestTasksCorrect(t *testing.T) {
	s := connectDb(t)

	// Создайте тестовую задачу
	task := storage.Task{
		ID:         3,
		Closed:     1524,
		AssignedID: 1,
		Title:      "Coorrect title",
		Content:    "Correct Test Content",
	}

	// Вызовите функцию TasksCorrect
	taskID, err := s.TasksCorrect(task)
	if err != nil {
		t.Fatalf("Failed to create new task: %v", err)
	}

	t.Log("Correct task with id:", taskID)

}

func TestDeleteTasks(t *testing.T) {
	s := connectDb(t)

	// Создайте тестовую задачу
	task := storage.Task{
		AuthorID: 3,
		Title:    "Delete Test Task",
		Content:  "Delete Test Content",
	}

	// Вызовите функцию NewTask
	taskID, err := s.NewTask(task)
	if err != nil {
		t.Fatalf("Failed to create new task: %v", err)
	}
	t.Log("Creat task with id:", taskID)

	err = s.DeleteTasks(taskID)
	if err != nil {
		t.Fatalf("Failed to delete new task: %v", err)
	}
	t.Log("Delete task with id:", taskID)

	tasks, err := s.GetTasks(taskID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}
	t.Log(tasks)

}
