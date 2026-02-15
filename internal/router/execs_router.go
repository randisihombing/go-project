package router

import (
	"gocourse/internal/handler"
	"net/http"
)

func execsRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /execs/", handler.GetExecsHandler)
	mux.HandleFunc("POST /execs/", handler.AddExecsHandler)
	mux.HandleFunc("PATCH /execs/", handler.PatchExecsHandler)

	mux.HandleFunc("GET /execs/{id}", handler.GetOneExecsHandler)
	mux.HandleFunc("PATCH /execs/{id}", handler.PatchOneExecsHandler)
	mux.HandleFunc("DELETE /execs/{id}", handler.DeleteOneExecsHandler)
	mux.HandleFunc("POST /execs/{id}/updatepassword", handler.GetExecsHandler)

	mux.HandleFunc("POST /execs/login", handler.LoginHandler)
	mux.HandleFunc("POST /execs/logout", handler.LogoutHandler)
	mux.HandleFunc("POST /execs/forgotpassword", handler.GetExecsHandler)
	mux.HandleFunc("POST /execs/resetpassword/reset/{resetcode}", handler.GetExecsHandler)

	return mux
}
