package config

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Route struct {
	Path   string
	Method string
}

type CustomMux struct {
	*http.ServeMux
	Routes   []Route
	NotFound http.HandlerFunc
}

func NewCustomMux() *CustomMux {
	return &CustomMux{
		ServeMux: http.NewServeMux(),
		NotFound: func(w http.ResponseWriter, r *http.Request) {
			AppLogger.Println("Route not found!")
			AppLogger.Println(r.URL.Host + r.URL.Path)

			response := map[string]string{
				"message": "Route Not Found: " + r.URL.Path,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
		},
	}
}

func (m *CustomMux) Handle(method, path string, Handler http.Handler) {
	m.Routes = append(m.Routes, Route{Method: method, Path: path})
	m.ServeMux.Handle(path, Handler)
}

func (m *CustomMux) HandleFunc(method, path string, Handler func(w http.ResponseWriter, r *http.Request)) {
	m.Routes = append(m.Routes, Route{Method: method, Path: path})
	m.ServeMux.HandleFunc(path, Handler)
}

func (m *CustomMux) ListRoutes() {
	AppLogger.Printf("Registered Routes:")
	for _, r := range m.Routes {
		AppLogger.Printf(strings.Join([]string{r.Method, r.Path}, " "))
	}
}
