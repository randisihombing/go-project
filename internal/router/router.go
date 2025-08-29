package router

import (
	"gocourse/internal/handler"
	"net/http"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.RootHandler)

	mux.HandleFunc("/teachers/", handler.TeachersHandler)

	mux.HandleFunc("/students/", handler.StudentsHandler)

	mux.HandleFunc("/execs/", handler.ExecsHandler)

	return mux
}
