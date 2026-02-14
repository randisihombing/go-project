package handler

import (
	"encoding/json"
	"fmt"
	"gocourse/internal/model"
	"gocourse/internal/repository/sqlconnect"
	"io"
	"log"
	"net/http"
	"strconv"
)

func GetExecsHandler(w http.ResponseWriter, r *http.Request) {

	var execs []model.Exec
	execs, err := sqlconnect.GetExecsDbHandler(execs, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := struct {
		Status string       `json:"status"`
		Count  int          `json:"count"`
		Data   []model.Exec `json:"data"`
	}{
		Status: "Success",
		Count:  len(execs),
		Data:   execs,
	}

	w.Header().Set("Content-Type", "applicaiton/json")
	json.NewEncoder(w).Encode(response)

}

// Get /exect/{id}
func GetOneExecsHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	exec, err := sqlconnect.GetExecByID(id)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exec)

}

func AddExecsHandler(w http.ResponseWriter, r *http.Request) {

	var newExecs []model.Exec
	var rawExects []map[string]interface{}

	body, error := io.ReadAll(r.Body)
	if error != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	err := json.Unmarshal(body, &rawExects)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	fields := GetFieldNames(model.Exec{})

	allowedFields := make(map[string]struct{})
	for _, field := range fields {
		allowedFields[field] = struct{}{}
	}

	for _, exec := range rawExects {
		for key := range exec {
			_, ok := allowedFields[key]
			if !ok {
				http.Error(w, "Unacceptable field found in request. Only use allowed fields.", http.StatusBadRequest)
				return
			}
		}
	}

	err = json.Unmarshal(body, &newExecs)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	for _, exec := range newExecs {
		err := CheckBlankFields(exec)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	addedExecs, err := sqlconnect.AddExecsDbHandler(newExecs)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string       `json:"status"`
		Count  int          `json:"count"`
		Data   []model.Exec `json:"data"`
	}{
		Status: "Success",
		Count:  len(addedExecs),
		Data:   addedExecs,
	}
	json.NewEncoder(w).Encode(response)
}

// Patch /exects/
func PatchExecsHandler(w http.ResponseWriter, r *http.Request) {

	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = sqlconnect.PatchExecs(updates)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

// Patch /execs/{id}
func PatchOneExecsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid exec Id", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}

	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request Id", http.StatusBadRequest)
		return
	}

	updatedExec, err := sqlconnect.PatchOneExec(id, updates)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedExec)
}

// Delete /execs/
func DeleteOneExecsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid exec Id", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteOneExec(id)
	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//response body
	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "Exec successfully deleted",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)
}
