package router

import (
	"gocourse/internal/handler"
	"net/http"
)

func execsRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /execs/", handler.ExecsHandler)
	mux.HandleFunc("POST /execs/", handler.ExecsHandler)
	mux.HandleFunc("PATCH /execs/", handler.ExecsHandler)

	mux.HandleFunc("GET /execs/{id}", handler.ExecsHandler)
	mux.HandleFunc("PATCH /execs/{id}", handler.ExecsHandler)
	mux.HandleFunc("DELETE /execs/{id}", handler.ExecsHandler)
	mux.HandleFunc("POST /execs/{id}/updatepassword", handler.ExecsHandler)

	mux.HandleFunc("POST /execs/login", handler.ExecsHandler)
	mux.HandleFunc("POST /execs/logout", handler.ExecsHandler)
	mux.HandleFunc("POST /execs/forgotpassword", handler.ExecsHandler)
	mux.HandleFunc("POST /execs/resetpassword/reset/{resetcode}", handler.ExecsHandler)

	return mux
}
