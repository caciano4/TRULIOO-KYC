package routes

import (
	"fmt"
	"net/http"
	"trullio-kyc/config"
	"trullio-kyc/controllers"
	"trullio-kyc/middleware"
	"trullio-kyc/utils"
)

func InitRoutes() {

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve main page
	http.HandleFunc("/", controllers.MainPage)

	// Get The list of PACKAGES and ADICIONAL INFORMATION
	http.Handle(
		"/kyc-package-list",
		utils.ChainMiddlewares(
			http.HandlerFunc(controllers.GetPackageList),
			middleware.CorsMiddleware,
			middleware.CheckMethodGet,
		),
	)

	// Submit KYC request
	http.Handle(
		"/kyc-request",
		utils.ChainMiddlewares(
			http.HandlerFunc(controllers.StoreFile),
			middleware.CorsMiddleware,
		),
	)

	// Process KYC Request
	http.Handle(
		"/process-kyc",
		utils.ChainMiddlewares(
			http.HandlerFunc(controllers.InitTrulioo),
			middleware.CorsMiddleware,
			middleware.CheckMethodGet,
		),
	)

	// TODO: Create route not found | Not exists

	StartServer()
}

func StartServer() {
	port := config.GetEnv("PORT", "80")
	fmt.Println("Server running on port " + port)
	http.ListenAndServe(":"+port, nil)
}
