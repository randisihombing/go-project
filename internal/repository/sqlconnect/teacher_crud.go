package sqlconnect

import (
	"database/sql"
	"fmt"
	"gocourse/internal/model"
	"gocourse/pkg/utils"
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

func GetTeachersDbHandler(teachers []model.Teacher, r *http.Request) ([]model.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {

		return nil, utils.ErrorHandler(err, "Error retrieving data")
	}

	defer db.Close()

	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"

	var args []interface{}

	query, args = addFilters(r, query, args)

	addSorting(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)

		return nil, utils.ErrorHandler(err, "Error retrieving data")
	}
	defer rows.Close()

	for rows.Next() {
		var teacher model.Teacher
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {

			return nil, utils.ErrorHandler(err, "Error retrieving data")
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

func GetTeacherByID(id int) (model.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {

		return model.Teacher{}, utils.ErrorHandler(err, "Error retrieving data")
	}

	defer db.Close()

	var teacher model.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
		&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {

		return model.Teacher{}, utils.ErrorHandler(err, "Error retrieving data")
	} else if err != nil {

		return model.Teacher{}, utils.ErrorHandler(err, "Error retrieving data")
	}
	return teacher, nil
}

func AddTeachersDbHandler(newTeachers []model.Teacher) ([]model.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {

		return nil, utils.ErrorHandler(err, "Error adding data")
	}

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	if err != nil {

		return nil, utils.ErrorHandler(err, "Error adding data")
	}
	defer stmt.Close()

	addedTeachers := make([]model.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {

			return nil, utils.ErrorHandler(err, "Error adding data")
		}
		lastID, err := res.LastInsertId()
		if err != nil {

			return nil, utils.ErrorHandler(err, "Error adding data")
		}
		newTeacher.ID = int(lastID)

		// teachers[nextID] = newTeacher
		addedTeachers[i] = newTeacher
		// nextID++
	}
	return addedTeachers, nil
}

func UpdateTeacher(id int, updatedTeacher model.Teacher) (model.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)

		return model.Teacher{}, utils.ErrorHandler(err, "Error update data")
	}
	defer db.Close()

	var existingTeacher model.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
		&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class,
		&existingTeacher.Subject)

	if err != nil {
		if err == sql.ErrNoRows {

			return model.Teacher{}, utils.ErrorHandler(err, "Error update data")
		}

		return model.Teacher{}, utils.ErrorHandler(err, "Error update data")
	}
	updatedTeacher.ID = existingTeacher.ID
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID)
	if err != nil {
		log.Println("update teacher error:", utils.ErrorHandler(err, "Error update data"))

		return model.Teacher{}, utils.ErrorHandler(err, "Error update data")
	}
	return updatedTeacher, nil
}

func PatchTeachers(updates []map[string]interface{}) error {
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)

		return utils.ErrorHandler(err, "Error update data")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {

		return utils.ErrorHandler(err, "Error update data")
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()

			return utils.ErrorHandler(err, "Invalid Id")
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()

			return utils.ErrorHandler(err, "Invalid Id")
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

				return utils.ErrorHandler(err, "Teacher not found")
			}

			return utils.ErrorHandler(err, "Error update data")
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
							return utils.ErrorHandler(err, "Error update data")
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

			return utils.ErrorHandler(err, "Error update data")
		}
	}
	err = tx.Commit()
	if err != nil {

		return utils.ErrorHandler(err, "Error update data")
	}
	return nil
}

func PatchOneTeacher(id int, updates map[string]interface{}) (model.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)

		return model.Teacher{}, utils.ErrorHandler(err, "Error update data")
	}
	defer db.Close()

	var existingTeacher model.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
		&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class,
		&existingTeacher.Subject)

	if err != nil {
		if err == sql.ErrNoRows {

			return model.Teacher{}, utils.ErrorHandler(err, "Teacher not found")
		}

		return model.Teacher{

			//Reflect package instead switch case
		}, utils.ErrorHandler(err, "Error update data")
	}

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

		return model.Teacher{}, utils.ErrorHandler(err, "Error update data")
	}
	return existingTeacher, nil
}

func DeleteOneTeacher(id int) error {
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)

		return utils.ErrorHandler(err, "Error update data")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {

		return utils.ErrorHandler(err, "Error update data")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {

		return utils.ErrorHandler(err, "Error update data")
	}

	if rowsAffected == 0 {

		return utils.ErrorHandler(err, "Teacher not found")
	}
	return nil
}

func DeleteTeachers(ids []int) ([]int, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)

		return nil, utils.ErrorHandler(err, "Error update data")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {

		return nil, utils.ErrorHandler(err, "Error update data")
	}

	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		tx.Rollback()

		return nil, utils.ErrorHandler(err, "Error update data")
	}
	defer stmt.Close()

	deletedIds := []int{}

	for _, id := range ids {
		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()

			return nil, utils.ErrorHandler(err, "Error update data")
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()

			return nil, utils.ErrorHandler(err, "Error update data")
		}
		//if teacher was deleted then add the id to the deleted ids slice
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}

		if rowsAffected < 1 {
			tx.Rollback()

			return nil, utils.ErrorHandler(err, fmt.Sprintf("ID %d not found", id))
		}

	}

	//commit
	err = tx.Commit()
	if err != nil {
		tx.Rollback()

		return nil, utils.ErrorHandler(err, "Error update data")
	}

	if len(deletedIds) < 1 {

		return nil, utils.ErrorHandler(err, "IDs do not exist")
	}
	return deletedIds, nil
}
