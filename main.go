package main

import (
	"github.com/gorilla/mux"
	"log"
	"myapp/config"
	"myapp/database"
	"myapp/handlers"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	// Настраиваем роутер
	r.HandleFunc("/", handlers.HomePage).Methods("GET")
	r.HandleFunc("/device", handlers.GetDeviceByMAC).Methods("POST")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Загружаем конфигурацию (если есть)
	config.LoadConfig()

	// Подключаемся к базе данных
	database.Connect()

	// Запуск сервера
	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
