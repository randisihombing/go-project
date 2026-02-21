package middlewares

import (
	"fmt"
	"net/http"
	"strings"
)

func MiddlewaresExcludePaths(middleware func(http.Handler) http.Handler, excludePaths ...string) func(http.Handler) http.Handler {
	fmt.Println("MiddlewareExcludePaths initialized")
	return func(next http.Handler) http.Handler {
		fmt.Println("====== MiddlewareExcludePaths RAN")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, path := range excludePaths {
				if strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
			}
			middleware(next).ServeHTTP(w, r)
			fmt.Println("Sent response from MiddlewareExcludePaths")
		})
	}
}
