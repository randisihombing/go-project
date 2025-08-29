package middlewares

import (
	"fmt"
	"net/http"
)

//api is hosted at www.myapi.com
//frontend server is at www.myfrontend.com

var allowedOrigin = []string{
	"https://my-origin-url.com",
	"https://www.myfrontend.com",
	"https://localhost:3000",
}

func Cors(next http.Handler) http.Handler {
	fmt.Println("Cors Middleware....")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Cors Middleware being returned....")
		origin := r.Header.Get("Origin")
		fmt.Println(origin)

		if isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			http.Error(w, "Not allowed by CORS", http.StatusForbidden)
			return
		}
		// w.Header().Set()
		w.Header().Set("Access-Controll-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Controll-Expose-Headers", "Authorization")
		w.Header().Set("Access-Controll-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Controll-Allow-Credetials", "true")
		w.Header().Set("Access-Controll-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
		fmt.Println("Cors Middleware ends....")
	})
}

func isOriginAllowed(origin string) bool {
	for _, alallowedOrigin := range allowedOrigin {
		if origin == alallowedOrigin {
			return true
		}
	}
	return false
}
