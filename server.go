package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "Hello Root Route")
	w.Write([]byte("Hello Root Route"))
	fmt.Println("Hello Root Router")
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {
	//teachers/{id}
	//teachers/?key=value&query=value2&sortby=email&sortorder=ASC
	fmt.Println(r.Method)
	switch r.Method {
	case http.MethodGet:
		fmt.Println(r.URL.Path)
		path := strings.TrimPrefix(r.URL.Path, "/teachers/")
		userID := strings.TrimSuffix(path, "/")

		fmt.Println("The ID is:", userID)
		fmt.Println("Query Params:", r.URL.Query())

		queryParams := r.URL.Query()
		sortBy := queryParams.Get("sortby")
		key := queryParams.Get("key")
		sortOrder := queryParams.Get("sortorder")

		if sortOrder == "" {
			sortOrder = "DESC"
		}

		fmt.Printf("Sort by: %v, Sort order: %v, Key: %v", sortBy, sortOrder, key)

		w.Write([]byte("Hello GET Method on Teachers Route"))
		// fmt.Println("Hello GET Method on Teachers Route")
	case http.MethodPost:
		w.Write([]byte("Hello POST Method on Teachers Route"))
		fmt.Println("Hello POST Method on Teachers Route")
	case http.MethodPut:
		w.Write([]byte("Hello PUT Method on Teachers Route"))
		fmt.Println("Hello PUT Method on Teachers Route")
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH Method on Teachers Route"))
		fmt.Println("Hello PATCH Method on Teachers Route")
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE Method on Teachers Route"))
		fmt.Println("Hello DELETE Method on Teachers Route")
	}
	// w.Write([]byte("Hello Teachers Route"))
	// fmt.Println("Hello Teachers Route")

}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET Method on Students Route"))
		fmt.Println("Hello GET Method on Students Route")
	case http.MethodPost:
		w.Write([]byte("Hello POST Method on Students Route"))
		fmt.Println("Hello POST Method on Students Route")
	case http.MethodPut:
		w.Write([]byte("Hello PUT Method on Students Route"))
		fmt.Println("Hello PUT Method on Students Route")
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH Method on Students Route"))
		fmt.Println("Hello PATCH Method on Students Route")
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE Method on Students Route"))
		fmt.Println("Hello DELETE Method on Students Route")
	}
	w.Write([]byte("Hello Students Route"))
	fmt.Println("Hello Students Route")
}

func execsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET Method on Execs Route"))
		fmt.Println("Hello GET Method on Execs Route")
	case http.MethodPost:
		w.Write([]byte("Hello POST Method on Execs Route"))
		fmt.Println("Hello POST Method on Execs Route")
	case http.MethodPut:
		w.Write([]byte("Hello PUT Method on Execs Route"))
		fmt.Println("Hello PUT Method on Execs Route")
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH Method on Execs Route"))
		fmt.Println("Hello PATCH Method on Execs Route")
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE Method on Execs Route"))
		fmt.Println("Hello DELETE Method on Execs Route")
	}
	w.Write([]byte("Hello Execs Route"))
	fmt.Println("Hello Execs Router")
}

func main() {
	port := ":3000"

	http.HandleFunc("/", rootHandler)

	http.HandleFunc("/teachers/", teachersHandler)

	http.HandleFunc("/students/", studentsHandler)

	http.HandleFunc("/execs/", execsHandler)

	fmt.Println("Server is running in port: ", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("Error starting the server", err)
	}
}
