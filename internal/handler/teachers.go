package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gocourse/internal/model"
	"gocourse/internal/repository/sqlconnect"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	return validFields[field]
}

func GetTeachersHandlers(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to DB", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	// path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	// idStr := strings.TrimSuffix(path, "/")
	// fmt.Println(idStr)

	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")

	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"

	var args []interface{}

	query, args = addFilters(r, query, args)

	addSorting(r, query)

	if firstName != "" {
		query += " AND first_name = ?"
		args = append(args, firstName)
	}
	if lastName != "" {
		query += " AND last_name = ?"
		args = append(args, lastName)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		http.Error(w, fmt.Sprintf("DB query error : %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	teacherList := make([]model.Teacher, 0)

	for rows.Next() {
		var teacher model.Teacher
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			http.Error(w, "Error scanning db results", http.StatusInternalServerError)
			return
		}
		teacherList = append(teacherList, teacher)
	}

	response := struct {
		Status string          `json:"status"`
		Count  int             `json:"count"`
		Data   []model.Teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(teacherList),
		Data:   teacherList,
	}

	w.Header().Set("Content-Type", "applicaiton/json")
	json.NewEncoder(w).Encode(response)

}

// Get /teachers/{id}
func GetOneTeacherHandlers(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to DB", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	idStr := r.PathValue("id")

	//handle path parameter
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	var teacher model.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
		&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "DB query error", http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)

}

func addSorting(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortby"]
	if len(sortParams) > 0 {
		query += " ORDER BY"
		for i, param := range sortParams {
			// sortBy=name:desc
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidSortField(field) || !isValidSortOrder(order) {
				continue
			}
			if i > 0 {
				query += ","
			}
			query += " " + field + " " + order
		}
	}
	return query
}

func addFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + " = ? "
			args = append(args, value)
		}

	}
	return query, args
}

func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to DB", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	var newTeachers []model.Teacher

	err = json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in preparing SQL query : %v", err), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	addedTeachers := make([]model.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			http.Error(w, "Error adding data into DB", http.StatusInternalServerError)
			return
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting last insert ID", http.StatusInternalServerError)
			return
		}
		newTeacher.ID = int(lastID)

		// teachers[nextID] = newTeacher
		addedTeachers[i] = newTeacher
		// nextID++
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

// Put /teachers/{id}
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid teacher Id", http.StatusBadRequest)
		return
	}

	var updatedTeacher model.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request Id", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable connect to DB", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var existingTeacher model.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
		&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class,
		&existingTeacher.Subject)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return
	}
	updatedTeacher.ID = existingTeacher.ID
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID)
	if err != nil {
		log.Println("update teacher error:", err)
		http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)
}

// Patch /teachers/
func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDb()
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable connect to DB", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var updates []map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			http.Error(w, "Invalid teacher ID in update", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error converting ID to int", http.StatusBadRequest)
			return
		}

		var teacherFromDb model.Teacher
		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
			&teacherFromDb.ID, &teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Email, &teacherFromDb.Class,
			&teacherFromDb.Subject)
		if err != nil {
			log.Println("ID:", id)
			log.Printf("Type: %T", id)
			log.Println(err)
			tx.Rollback()
			if err == sql.ErrNoRows {
				http.Error(w, "Teacher not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Error retrieving teacher", http.StatusInternalServerError)
			return
		}

		//Apply update using reflection
		teacherVal := reflect.ValueOf(&teacherFromDb).Elem()
		teacherType := teacherVal.Type()

		for k, v := range update {
			if k == "id" {
				//skip updating id field
				continue
			}
			for i := 0; i < teacherVal.NumField(); i++ {
				field := teacherType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := teacherVal.Field(i)
					if teacherVal.Field(i).CanSet() {
						val := reflect.ValueOf(v)
						// fieldVal.Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							tx.Rollback()
							log.Printf("Cannot convert %v to %v", val.Type(), fieldVal.Type())
							return
						}
					}
					break
				}
			}
		}

		_, err = tx.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
			teacherFromDb.FirstName, teacherFromDb.LastName, teacherFromDb.Email, teacherFromDb.Class, teacherFromDb.Subject, teacherFromDb.ID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error updating teacher", http.StatusInternalServerError)
			return
		}

		err = tx.Commit()
		if err != nil {
			http.Error(w, "Error committing transaction", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(teacherFromDb)
	}

}

// Patch /teachers/{id}
func PatchOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid teacher Id", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}

	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request Id", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable connect to DB", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var existingTeacher model.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
		&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class,
		&existingTeacher.Subject)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return
	}

	//Reflect package instead switch case
	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	teacherType := teacherVal.Type()

	for k, v := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			field.Tag.Get("json")
			if field.Tag.Get("json") == k+" ,omitempty" {
				if teacherVal.Field(i).CanSet() {
					fieldVal := teacherVal.Field(i)
					fieldVal.Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID)
	if err != nil {
		log.Println("update teacher error:", err)
		http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingTeacher)
}

// Delete /teachers/
func DeleteOneTeacherhandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid teacher Id", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable connect to DB", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Error deleting teacher", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error retrieving result", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	// w.WriteHeader(http.StatusNoContent)
	// json.NewEncoder(w).Encode(existingTeacher)

	//response body
	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "Teacher successfully deleted",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDb()
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable connect to DB", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var ids []int
	err = json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}

	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error Preparing delete statemend", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	// if err != nil {
	// 	http.Error(w, "Error deleting teacher", http.StatusInternalServerError)
	// 	return
	// }

	deletedIds := []int{}

	for _, id := range ids {
		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error deleting teacher", http.StatusInternalServerError)
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error retrieving result", http.StatusInternalServerError)
			return
		}
		//if teacher was deleted then add the id to the deleted ids slice
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}

		if rowsAffected < 1 {
			tx.Rollback()
			http.Error(w, fmt.Sprintf("ID %ddoes not exist", id), http.StatusInternalServerError)
			return
		}

	}

	//commit
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error Committing transaction", http.StatusInternalServerError)
		return
	}

	if len(deletedIds) < 1 {
		http.Error(w, "IDs do not exist", http.StatusInternalServerError)
		return
	}
	//response body
	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status     string `json:"status"`
		DeletedIds []int  `json:"id"`
	}{
		Status:     "Teacher successfully deleted",
		DeletedIds: deletedIds,
	}
	json.NewEncoder(w).Encode(response)

}
