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

func GetStudentsHandlers(w http.ResponseWriter, r *http.Request) {

	var students []model.Student
	students, err := sqlconnect.GetStudentsDbHandler(students, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := struct {
		Status string          `json:"status"`
		Count  int             `json:"count"`
		Data   []model.Student `json:"data"`
	}{
		Status: "Success",
		Count:  len(students),
		Data:   students,
	}

	w.Header().Set("Content-Type", "applicaiton/json")
	json.NewEncoder(w).Encode(response)

}

// Get /students/{id}
func GetOneStudentHandlers(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	student, err := sqlconnect.GetStudentByID(id)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)

}

func AddStudentHandler(w http.ResponseWriter, r *http.Request) {

	var newStudents []model.Student
	var rawStudents []map[string]interface{}

	body, error := io.ReadAll(r.Body)
	if error != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	err := json.Unmarshal(body, &rawStudents)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	fields := GetFieldNames(model.Student{})

	allowedFields := make(map[string]struct{})
	for _, field := range fields {
		allowedFields[field] = struct{}{}
	}

	for _, student := range rawStudents {
		for key := range student {
			_, ok := allowedFields[key]
			if !ok {
				http.Error(w, "Unacceptable field found in request. Only use allowed fields.", http.StatusBadRequest)
				return
			}
		}
	}

	err = json.Unmarshal(body, &newStudents)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	for _, student := range newStudents {
		err := CheckBlankFields(student)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	addedStudents, err := sqlconnect.AddStudentsDbHandler(newStudents)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string          `json:"status"`
		Count  int             `json:"count"`
		Data   []model.Student `json:"data"`
	}{
		Status: "Success",
		Count:  len(addedStudents),
		Data:   addedStudents,
	}
	json.NewEncoder(w).Encode(response)
}

// Put /students/{id}
func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid student Id", http.StatusBadRequest)
		return
	}

	var updatedStudent model.Student
	err = json.NewDecoder(r.Body).Decode(&updatedStudent)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request Id", http.StatusBadRequest)
		return
	}

	updatedStudentFromDB, err := sqlconnect.UpdateStudent(id, updatedStudent)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedStudentFromDB)
}

// Patch /students/
func PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {

	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = sqlconnect.PatchStudents(updates)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

// Patch /Students/{id}
func PatchOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid student Id", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}

	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request Id", http.StatusBadRequest)
		return
	}

	updatedStudent, err := sqlconnect.PatchOneStudent(id, updates)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedStudent)
}

// Delete /students/
func DeleteOneStudenthandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid student Id", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteOneStudent(id)
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
		Status: "Student successfully deleted",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteStudentsHandler(w http.ResponseWriter, r *http.Request) {

	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Print(err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	deletedIds, err := sqlconnect.DeleteStudents(ids)
	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//response body
	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status     string `json:"status"`
		DeletedIds []int  `json:"id"`
	}{
		Status:     "Student successfully deleted",
		DeletedIds: deletedIds,
	}
	json.NewEncoder(w).Encode(response)

}
