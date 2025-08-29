package middlewares

import (
	"fmt"
	"net/http"
)

func SecurityHeaders(next http.Handler) http.Handler {
	fmt.Println("Security Headers Middleware....")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Security Headers Middleware being returned....")
		w.Header().Set("X-DNS-Prefetch-Control", "off")
		w.Header().Set("X-Frame-Options", "off")
		w.Header().Set("X-XSS-Protection", "1;mode=block")
		w.Header().Set("X-Content-Type-Options", "nonsniff")
		w.Header().Set("Strict Transsport Security", "max-age=63072000;includeSubDomains;preload")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Reffered-Policy", "no-referrer")
		w.Header().Set("X-Powered-By", "Django")
		w.Header().Set("Server", "")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-aged=0")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Permission-Policy", "geolocation=(self), microphone=()")
		next.ServeHTTP(w, r)
		fmt.Println("Security Headers Middleware ends....")

	})
}

//Basic Middleware Skeleton
// func securityHeaders(next http.Handler) http.Handler{
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		next.ServeHTTP(w, r)
// 	})
// }
