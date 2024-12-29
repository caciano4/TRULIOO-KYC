package routes

import (
	"net/http"
	"trullio-kyc/config"
	"trullio-kyc/controllers"
	"trullio-kyc/middleware"
	"trullio-kyc/utils"
)

func InitRoutes() {
	mux := config.NewCustomMux()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET", "/static/", http.StripPrefix("/static/", fs))

	// Serve main page
	mux.HandleFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			controllers.MainPage(w, r)
		} else {
			mux.NotFound(w, r)
		}
	})

	// Get The list of PACKAGES and ADICIONAL INFORMATION
	mux.Handle(
		"GET",
		"/kyc-package-list",
		utils.ChainMiddlewares(
			http.HandlerFunc(controllers.GetPackageList),
			middleware.CorsMiddleware,
			middleware.CheckMethodGet,
		),
	)

	// Submit KYC request
	mux.Handle(
		"POST",
		"/kyc-request",
		utils.ChainMiddlewares(
			http.HandlerFunc(controllers.StoreFile),
			middleware.CheckMethodPost,
			middleware.CorsMiddleware,
		),
	)

	// Process KYC Request
	mux.Handle(
		"GET",
		"/process-kyc/",
		utils.ChainMiddlewares(
			http.HandlerFunc(controllers.InitTrulioo),
			middleware.CorsMiddleware,
			middleware.CheckMethodGet,
			middleware.ExtractParamMiddleware,
		),
	)

	StartServer(mux)
}

func StartServer(mux *config.CustomMux) {
	port := config.GetEnv("PORT", "80")
	mux.ListRoutes()
	config.AppLogger.Println("Server running on port " + port)
	http.ListenAndServe(":"+port, mux)
}
