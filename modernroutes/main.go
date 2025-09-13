package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	//Method based routing
	mux.HandleFunc("POST /items/create", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Item created")
	})

	mux.HandleFunc("DELETE /items/create", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Item deleted")
	})

	//Wildcard in pattern - path pattern
	mux.HandleFunc("GET /teachers/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Teacher ID: %s", r.PathValue("id"))
	})

	//Wildcard with ...
	mux.HandleFunc("/files/{path...}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Path: %s", r.PathValue("path"))
	})

	mux.HandleFunc("/path1/{param1}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Param1: %s", r.PathValue("param1"))
	})

	// mux.HandleFunc("/{param2}/path2", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Param2: %s", r.PathValue("param2"))
	// })

	mux.HandleFunc("/path1/path2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Nothing to see here")
	})

	http.ListenAndServe(":8080", mux)
}
