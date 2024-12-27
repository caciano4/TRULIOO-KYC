package middleware

import "net/http"

func CheckMethodGet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{message: "Method not Allowrd"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CheckMethodPost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CheckMethodDelete(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{message: "Method not Allowed"}`))
		}
	})
}
