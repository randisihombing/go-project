package router

import (
	"gocourse/internal/handler"
	"net/http"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handler.RootHandler)

	//Teacher
	mux.HandleFunc("GET /teachers/", handler.GetTeachersHandlers)
	mux.HandleFunc("POST /teachers/", handler.AddTeacherHandler)
	mux.HandleFunc("PUT /teachers/", handler.UpdateTeacherHandler)
	mux.HandleFunc("PATCH /teachers/", handler.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teachers/", handler.DeleteTeachersHandler)

	mux.HandleFunc("PUT /teachers/{id}", handler.UpdateTeacherHandler)
	mux.HandleFunc("GET /teachers/{id}", handler.GetOneTeacherHandlers)
	mux.HandleFunc("PATCH /teachers/{id}", handler.PatchOneTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handler.DeleteOneTeacherhandler)

	mux.HandleFunc("GET /teachers/{id}/students", handler.GetStudentsByTeacherId)
	mux.HandleFunc("GET /teachers/{id}/studentcount", handler.GetTeachersHandlers)

	//Student
	mux.HandleFunc("GET /students/", handler.GetStudentsHandlers)
	mux.HandleFunc("POST /students/", handler.AddStudentHandler)
	mux.HandleFunc("PUT /students/", handler.UpdateStudentHandler)
	mux.HandleFunc("PATCH /students/", handler.PatchStudentsHandler)
	mux.HandleFunc("DELETE /students/", handler.DeleteStudentsHandler)

	mux.HandleFunc("PUT /students/{id}", handler.UpdateStudentHandler)
	mux.HandleFunc("GET /students/{id}", handler.GetOneStudentHandlers)
	mux.HandleFunc("PATCH /students/{id}", handler.PatchOneStudentHandler)
	mux.HandleFunc("DELETE /students/{id}", handler.DeleteOneStudenthandler)

	mux.HandleFunc("GET /execs/", handler.ExecsHandler)

	return mux
}
