package routes

import (
	"fmt"
	"net/http"
	"trullio-kyc/config"
	"trullio-kyc/controllers"
	"trullio-kyc/middleware"
	"trullio-kyc/utils"
)

// type Route struct {
// 	Method string
// 	Path   string
// }

// type CustomMux struct {
// 	*http.ServeMux
// 	Routes []Route
// }

// func newCustomMux() *CustomMux {
// 	return &CustomMux{
// 		ServeMux: http.NewServeMux(),
// 	}
// }

// func (*CustomMux)

func InitRoutes() {

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve main page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			controllers.MainPage(w, r)
		} else {
			config.AppLogger.Print("deu Merda!")
			http.NotFound(w, r)
		}
	})

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
			middleware.ExtractParamMiddleware(),
		),
	)

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	config.AppLogger.Print("teste")
	// 	w.WriteHeader(http.StatusNotFound)
	// 	response := map[string]interface{}{
	// 		"message": "Rsoute not found.",
	// 	}

	// 	json.NewEncoder(w).Encode(response)
	// })

	StartServer()
}

func StartServer() {
	port := config.GetEnv("PORT", "80")
	fmt.Println("Server running on port " + port)
	http.ListenAndServe(":"+port, nil)
}
