package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go_project/app"
	"go_project/config/data_base"
)

func main() {
	// Подключение к базе данных
	configPath := "./config/data_base/DBConnect.json"
	// configPath := "./config/data_base/TestDBConnect.json"
	db, err := data_base.OpenConnection(configPath)
	if err != nil {
		log.Fatalf("Couldn't connect to the database: %v", err)
	}
	defer db.Close()

	// Миграция схемы базы данных (создаем/актуализируем структуру данных)
	data_base.Migrate(db)

	// Создание маршрутизатора
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tasks", app.CreateTaskHandler(db)).Methods("POST")
	router.HandleFunc("/mytasks", app.GetTasksHandler(db)).Methods("GET")
	router.HandleFunc("/tasks/{id}", app.GetTaskHandler(db)).Methods("GET")
	router.HandleFunc("/tasks/{id}", app.CompleteTaskHandler(db)).Methods("PUT")
	router.HandleFunc("/tasks/{id}", app.DeleteTaskHandler(db)).Methods("DELETE")
	router.HandleFunc("/register", app.CreateUser(db)).Methods("POST")
	router.Handle("/metrics", promhttp.Handler())

	fmt.Println("The server is running on the port :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
