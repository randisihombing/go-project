package router

import (
	"gocourse/internal/handler"
	"net/http"
)

func studentsRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /students/", handler.GetStudentsHandlers)
	mux.HandleFunc("POST /students/", handler.AddStudentHandler)
	mux.HandleFunc("PUT /students/", handler.UpdateStudentHandler)
	mux.HandleFunc("PATCH /students/", handler.PatchStudentsHandler)
	mux.HandleFunc("DELETE /students/", handler.DeleteStudentsHandler)

	mux.HandleFunc("PUT /students/{id}", handler.UpdateStudentHandler)
	mux.HandleFunc("GET /students/{id}", handler.GetOneStudentHandlers)
	mux.HandleFunc("PATCH /students/{id}", handler.PatchOneStudentHandler)
	mux.HandleFunc("DELETE /students/{id}", handler.DeleteOneStudenthandler)

	return mux
}
