package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"trullio-kyc/config"
	"trullio-kyc/controllers"
)

func InitRoutes(db *sql.DB) {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve main page
	http.HandleFunc("/", controllers.MainPage)

	// Submit KYC request
	http.HandleFunc("/kyc-request", func(w http.ResponseWriter, r *http.Request) {
		controllers.StoreFile(db, w, r)
	})

	StartServer()
}

func StartServer() {
	port := config.GetEnv("PORT", "80")
	fmt.Println("Server running on port " + port)
	http.ListenAndServe(":"+port, nil)
}
