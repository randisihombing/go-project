package router

import (
	"gocourse/internal/handler"
	"net/http"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.RootHandler)

	mux.HandleFunc("GET /teachers/", handler.GetTeachersHandlers)
	mux.HandleFunc("POST /teachers/", handler.AddTeacherHandler)
	mux.HandleFunc("PUT /teachers/", handler.UpdateTeacherHandler)
	mux.HandleFunc("PATCH /teachers/", handler.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teachers/", handler.DeleteTeachersHandler)

	mux.HandleFunc("PUT /teachers/{id}", handler.UpdateTeacherHandler)
	mux.HandleFunc("GET /teachers/{id}", handler.GetOneTeacherHandlers)
	mux.HandleFunc("PATCH /teachers/{id}", handler.PatchOneTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handler.DeleteOneTeacherhandler)

	mux.HandleFunc("/students/", handler.StudentsHandler)

	mux.HandleFunc("/execs/", handler.ExecsHandler)

	return mux
}
