package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go_project/app"
	"go_project/config/data_base"
)

// Функция для тестирования запросов к API
func TestAPI(t *testing.T) {
	// t.Parallel()

	// Создаём фиктивную базу данных для тестирования
	mockDB, err := setupMockDatabase()
	require.NoError(t, err)
	defer mockDB.Close()

	truncateTable(mockDB)

	// Проверяем создание новой задачи
	t.Run("CreateTaskHandler", func(t *testing.T) {
		truncateTable(mockDB)
		newTask := map[string]interface{}{
			"title":       "CreateTaskHandler Test 1",
			"description": "Description Test 1",
		}

		reqBody, err := json.Marshal(newTask)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		handler := app.CreateTaskHandler(mockDB)
		handler.ServeHTTP(recorder, req)

		resp := recorder.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var createdTask data_base.Task
		err = json.NewDecoder(resp.Body).Decode(&createdTask)
		require.NoError(t, err)

		assert.NotZero(t, createdTask.ID)
		assert.Equal(t, newTask["title"], createdTask.Title)
		assert.Equal(t, newTask["description"], createdTask.Description)
	})

	// Проверяем получение списка всех задач
	t.Run("GetTasksHandler", func(t *testing.T) {
		truncateTable(mockDB)
		newTask := map[string]interface{}{
			"title":       "CreateTaskHandler Test 2",
			"description": "Description Test 2",
		}

		reqBody, err := json.Marshal(newTask)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		handler := app.CreateTaskHandler(mockDB)
		handler.ServeHTTP(recorder, req)

		resp := recorder.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var createdTask data_base.Task
		err = json.NewDecoder(resp.Body).Decode(&createdTask)
		require.NoError(t, err)

		assert.NotZero(t, createdTask.ID)

		req, err = http.NewRequest("GET", "/tasks", nil)
		require.NoError(t, err)

		recorder = httptest.NewRecorder()
		handler = app.GetTasksHandler(mockDB)
		handler.ServeHTTP(recorder, req)

		resp = recorder.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var tasks []data_base.Task
		err = json.NewDecoder(resp.Body).Decode(&tasks)
		require.NoError(t, err)

		assert.Len(t, tasks, 1)
	})

	// Проверяем получение одной конкретной задачи
	t.Run("GetTaskHandler", func(t *testing.T) {
		truncateTable(mockDB)
		newTask := map[string]interface{}{
			"title":       "CreateTaskHandler Test 2",
			"description": "Description Test 2",
		}

		reqBody, err := json.Marshal(newTask)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		handler := app.CreateTaskHandler(mockDB)
		handler.ServeHTTP(recorder, req)

		resp := recorder.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var createdTask data_base.Task
		err = json.NewDecoder(resp.Body).Decode(&createdTask)
		require.NoError(t, err)

		assert.NotZero(t, createdTask.ID)

		// Получаем первую задачу из тестового набора
		task := createdTask

		req, err = http.NewRequest("GET", fmt.Sprintf("/tasks/%d", task.ID), nil)
		require.NoError(t, err)

		recorder = httptest.NewRecorder()
		handler = app.GetTaskHandler(mockDB)
		handler.ServeHTTP(recorder, req)

		resp = recorder.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var fetchedTask data_base.Task
		err = json.NewDecoder(resp.Body).Decode(&fetchedTask)
		require.NoError(t, err)

		assert.Equal(t, task.ID, fetchedTask.ID)
		assert.Equal(t, task.Title, fetchedTask.Title)
		assert.Equal(t, task.Description, fetchedTask.Description)
		assert.Equal(t, task.IsCompleted, fetchedTask.IsCompleted)
	})

	// Проверяем обновление существующей задачи
	t.Run("CompleteTaskHandler", func(t *testing.T) {
		truncateTable(mockDB)
		newTask := map[string]interface{}{
			"title":       "CreateTaskHandler Test 2",
			"description": "Description Test 2",
		}

		reqBody, err := json.Marshal(newTask)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		handler := app.CreateTaskHandler(mockDB)
		handler.ServeHTTP(recorder, req)

		resp := recorder.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var createdTask data_base.Task
		err = json.NewDecoder(resp.Body).Decode(&createdTask)
		require.NoError(t, err)

		assert.NotZero(t, createdTask.ID)

		// Получаем первую задачу из тестового набора
		task := createdTask

		updatedTask := data_base.Task{
			Title:       "Updated Test Task",
			Description: "This is an updated test task",
			IsCompleted: true,
		}

		reqBody, err = json.Marshal(updatedTask)
		require.NoError(t, err)

		req, err = http.NewRequest("PUT", fmt.Sprintf("/tasks/%d", task.ID), bytes.NewBuffer(reqBody))
		require.NoError(t, err)

		recorder = httptest.NewRecorder()
		handler = app.CompleteTaskHandler(mockDB)
		handler.ServeHTTP(recorder, req)

		resp = recorder.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var completedTask data_base.Task
		err = json.NewDecoder(resp.Body).Decode(&completedTask)
		require.NoError(t, err)

		assert.Equal(t, updatedTask.Title, completedTask.Title)
		assert.Equal(t, updatedTask.Description, completedTask.Description)
		assert.True(t, completedTask.IsCompleted)
	})

	// Проверяем удаление задачи
	t.Run("DeleteTaskHandler", func(t *testing.T) {
		truncateTable(mockDB)
		newTask := map[string]interface{}{
			"title":       "CreateTaskHandler Test 2",
			"description": "Description Test 2",
		}

		reqBody, err := json.Marshal(newTask)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		handler := app.CreateTaskHandler(mockDB)
		handler.ServeHTTP(recorder, req)

		resp := recorder.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var createdTask data_base.Task
		err = json.NewDecoder(resp.Body).Decode(&createdTask)
		require.NoError(t, err)

		assert.NotZero(t, createdTask.ID)

		// Получаем первую задачу из тестового набора
		task := createdTask

		req, err = http.NewRequest("DELETE", fmt.Sprintf("/tasks/%d", task.ID), nil)
		require.NoError(t, err)

		recorder = httptest.NewRecorder()
		handler = app.DeleteTaskHandler(mockDB)
		handler.ServeHTTP(recorder, req)

		resp = recorder.Result()
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Проверяем, что задача была удалена
		var deletedTask data_base.Task
		err = mockDB.Where("id = ?", task.ID).First(&deletedTask).Error
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

// Настройка фиктивной базы данных для тестирования
func setupMockDatabase() (*gorm.DB, error) {
	configPath := "../config/data_base/TestDBConnect.json"
	mockDB, err := data_base.OpenConnection(configPath)
	if err != nil {
		return nil, err
	}

	// Выполняем миграцию структуры данных
	mockDB.AutoMigrate(&data_base.Task{})

	return mockDB, nil
}

// Функция для очистки таблицы перед каждым тестом
func truncateTable(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE tasks;")
}
