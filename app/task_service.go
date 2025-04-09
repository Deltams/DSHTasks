package app

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gopkg.in/go-playground/validator.v9"

	"go_project/config/data_base"
)

var validate *validator.Validate
var trans ut.Translator

var (
	activeUsers = promauto.NewCounter(prometheus.CounterOpts{
		Name: "task_service_processed_active_users",
		Help: "Total registered users since the start of the system",
	})
)

var (
	createTasks = promauto.NewCounter(prometheus.CounterOpts{
		Name: "task_service_processed_create_tasks",
		Help: "Total tasks created since system launch",
	})
)

func init() {
	validate = validator.New()
	uni := ut.New(en.New())
	trans, _ = uni.GetTranslator("en")
}

func CreateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var myUser data_base.MyUser

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading the request body", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &myUser)
		if err != nil {
			http.Error(w, "JSON parsing error", http.StatusBadRequest)
			return
		}

		result := db.Create(&myUser)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		activeUsers.Inc()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(myUser)
	}
}

func CreateTaskHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task data_base.Task

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading the request body", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &task)
		if err != nil {
			http.Error(w, "JSON parsing error", http.StatusBadRequest)
			return
		}

		err = validate.Struct(task)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				http.Error(w, err.Translate(trans), http.StatusUnprocessableEntity)
				return
			}
		}

		result := db.Create(&task)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		createTasks.Inc()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}
}

func GetTaskHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var task data_base.Task
		result := db.First(&task, id)
		if result.RecordNotFound() {
			http.NotFound(w, r)
			return
		}

		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}
}

func GetTasksHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tasks []data_base.Task
		result := db.Find(&tasks)
		if result.RecordNotFound() {
			http.NotFound(w, r)
			return
		}

		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}
}

func CompleteTaskHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task_db data_base.Task
		var task_new data_base.Task

		params := mux.Vars(r)
		id := params["id"]

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading the request body", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &task_new)
		if err != nil {
			http.Error(w, "JSON parsing error", http.StatusBadRequest)
			return
		}

		err = validate.Struct(task_new)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				http.Error(w, err.Translate(trans), http.StatusUnprocessableEntity)
				return
			}
		}

		result := db.First(&task_db, id)
		if result.RecordNotFound() {
			http.NotFound(w, r)
			return
		}

		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		task_db.Title = task_new.Title
		task_db.Description = task_new.Description
		task_db.IsCompleted = task_new.IsCompleted
		db.Save(&task_db)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task_db)
	}
}

func DeleteTaskHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var task data_base.Task
		result := db.First(&task, id)
		if result.RecordNotFound() {
			http.NotFound(w, r)
			return
		}

		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		db.Delete(&task)

		w.WriteHeader(http.StatusNoContent)
	}
}
