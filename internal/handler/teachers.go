package handler

import (
	"encoding/json"
	"fmt"
	"gocourse/internal/model"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	teachers = make(map[int]model.Teacher)
	mutex    = &sync.Mutex{}
	nextID   = 1
)

// initialize some dummy data
func init() {
	teachers[nextID] = model.Teacher{
		ID:        nextID,
		FirstName: "Jhon",
		LasttName: "Doe",
		Class:     "9A",
		Subject:   "Math",
	}
	nextID++
	teachers[nextID] = model.Teacher{
		ID:        nextID,
		FirstName: "Jane",
		LasttName: "Smith",
		Class:     "10A",
		Subject:   "Algebra",
	}
	nextID++
	teachers[nextID] = model.Teacher{
		ID:        nextID,
		FirstName: "Jane",
		LasttName: "Doe",
		Class:     "11A",
		Subject:   "Biology",
	}
	nextID++
}

func getTeachersHandlers(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")
	fmt.Println(idStr)

	if idStr == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")

		teacherList := make([]model.Teacher, 0, len(teachers))

		for _, teacher := range teachers {
			if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LasttName == lastName) {
				teacherList = append(teacherList, teacher)
			}
		}
		response := struct {
			Status string          `json:"status"`
			Count  int             `json:"count"`
			Data   []model.Teacher `json:"data"`
		}{
			Status: "Success",
			Count:  len(teachers),
			Data:   teacherList,
		}

		w.Header().Set("Content-Type", "applicaiton/json")
		json.NewEncoder(w).Encode(response)
	}

	//handle path parameter
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	teacher, exists := teachers[id]
	if !exists {
		http.Error(w, "Teachers not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(teacher)

}

func addTeacherHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var newTeachers []model.Teacher

	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	addedTeachers := make([]model.Teacher, len(newTeachers))

	for i, newTeacher := range newTeachers {
		newTeacher.ID = nextID
		teachers[nextID] = newTeacher
		addedTeachers[i] = newTeacher
		nextID++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string          `json:"status"`
		Count  int             `json:"count"`
		Data   []model.Teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	json.NewEncoder(w).Encode(response)
}

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	switch r.Method {
	case http.MethodGet:
		//call get method handler function
		getTeachersHandlers(w, r)
	case http.MethodPost:
		addTeacherHandler(w, r)
	case http.MethodPut:
		w.Write([]byte("Hello PUT Method on Teachers Route"))
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH Method on Teachers Route"))
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE Method on Teachers Route"))
	}
}
