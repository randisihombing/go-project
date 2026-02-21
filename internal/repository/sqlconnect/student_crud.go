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

func GetStudentsDbHandler(students []model.Student, r *http.Request, limit, page int) ([]model.Student, int, error) {
	db, err := ConnectDb()
	if err != nil {

		return nil, 0, utils.ErrorHandler(err, "Error retrieving data")
	}

	defer db.Close()

	query := "SELECT id, first_name, last_name, email, class FROM students WHERE 1=1"
	// LIMIT = ? OFFSET = ?

	var args []interface{}

	query, args = utils.AddFilters(r, query, args)

	//Add pagination
	offset := (page - 1) * limit
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	utils.AddSorting(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)

		return nil, 0, utils.ErrorHandler(err, "Error retrieving data")
	}
	defer rows.Close()

	for rows.Next() {
		var student model.Student
		err := rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {

			return nil, 0, utils.ErrorHandler(err, "Error retrieving data")
		}
		students = append(students, student)
	}

	//Get total count
	var totalStudents int
	err = db.QueryRow("SELECT COUNT(*) FROM students").Scan(&totalStudents)
	if err != nil {
		utils.ErrorHandler(err, "")
		totalStudents = 0
	}
	return students, totalStudents, nil
}

func GetStudentByID(id int) (model.Student, error) {
	db, err := ConnectDb()
	if err != nil {

		return model.Student{}, utils.ErrorHandler(err, "Error retrieving data")
	}

	defer db.Close()

	var student model.Student
	err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(
		&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
	if err == sql.ErrNoRows {

		return model.Student{}, utils.ErrorHandler(err, "Error retrieving data")
	} else if err != nil {

		return model.Student{}, utils.ErrorHandler(err, "Error retrieving data")
	}
	return student, nil
}

func AddStudentsDbHandler(newStudents []model.Student) ([]model.Student, error) {
	db, err := ConnectDb()
	if err != nil {

		return nil, utils.ErrorHandler(err, "Error adding data")
	}

	defer db.Close()

	stmt, err := db.Prepare(utils.GenerateInsertQuery("students", model.Student{}))
	if err != nil {

		return nil, utils.ErrorHandler(err, "Error adding data")
	}
	defer stmt.Close()

	addedStudents := make([]model.Student, len(newStudents))
	for i, newStudent := range newStudents {
		values := utils.GetStructValues(newStudent)
		res, err := stmt.Exec(values...)
		if err != nil {
			fmt.Println("----Error:", err.Error())
			if strings.Contains(err.Error(), "a foreign key constraint fails (`school`.`students`, CONSTRAINT `students_ibfk_1` FOREIGN KEY (`class`) REFERENCES`teachers` (`class`))") {
				return nil, utils.ErrorHandler(err, "class/class teacher does not exist")
			}
			return nil, utils.ErrorHandler(err, "Error adding data")
		}
		lastID, err := res.LastInsertId()
		if err != nil {

			return nil, utils.ErrorHandler(err, "Error adding data")
		}
		newStudent.ID = int(lastID)

		addedStudents[i] = newStudent
		// nextID++
	}
	return addedStudents, nil
}

func UpdateStudent(id int, updatedStudent model.Student) (model.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)

		return model.Student{}, utils.ErrorHandler(err, "Error update data")
	}
	defer db.Close()

	var existingStudent model.Student
	err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(
		&existingStudent.ID, &existingStudent.FirstName, &existingStudent.LastName, &existingStudent.Email, &existingStudent.Class)

	if err != nil {
		if err == sql.ErrNoRows {

			return model.Student{}, utils.ErrorHandler(err, "Error update data")
		}

		return model.Student{}, utils.ErrorHandler(err, "Error update data")
	}
	updatedStudent.ID = existingStudent.ID
	_, err = db.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?",
		updatedStudent.FirstName, updatedStudent.LastName, updatedStudent.Email, updatedStudent.Class, updatedStudent.ID)
	if err != nil {
		log.Println("update student error:", utils.ErrorHandler(err, "Error update data"))

		return model.Student{}, utils.ErrorHandler(err, "Error update data")
	}
	return updatedStudent, nil
}

func PatchStudents(updates []map[string]interface{}) error {
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

		var studentFromDb model.Student
		err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(
			&studentFromDb.ID, &studentFromDb.FirstName, &studentFromDb.LastName, &studentFromDb.Email, &studentFromDb.Class)
		if err != nil {
			log.Println("ID:", id)
			log.Printf("Type: %T", id)
			log.Println(err)
			tx.Rollback()
			if err == sql.ErrNoRows {

				return utils.ErrorHandler(err, "Student not found")
			}

			return utils.ErrorHandler(err, "Error update data")
		}

		//Apply update using reflection
		studentVal := reflect.ValueOf(&studentFromDb).Elem()
		studentType := studentVal.Type()

		for k, v := range update {
			if k == "id" {
				//skip updating id field
				continue
			}
			for i := 0; i < studentVal.NumField(); i++ {
				field := studentType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := studentVal.Field(i)
					if studentVal.Field(i).CanSet() {
						val := reflect.ValueOf(v)
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

		_, err = tx.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?",
			studentFromDb.FirstName, studentFromDb.LastName, studentFromDb.Email, studentFromDb.Class, studentFromDb.ID)
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

func PatchOneStudent(id int, updates map[string]interface{}) (model.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)

		return model.Student{}, utils.ErrorHandler(err, "Error update data")
	}
	defer db.Close()

	var existingStudent model.Student
	err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(
		&existingStudent.ID, &existingStudent.FirstName, &existingStudent.LastName, &existingStudent.Email, &existingStudent.Class)

	if err != nil {
		if err == sql.ErrNoRows {

			return model.Student{}, utils.ErrorHandler(err, "Student not found")
		}

		return model.Student{

			//Reflect package instead switch case
		}, utils.ErrorHandler(err, "Error update data")
	}

	studentVal := reflect.ValueOf(&existingStudent).Elem()
	studentType := studentVal.Type()

	for k, v := range updates {
		for i := 0; i < studentVal.NumField(); i++ {
			field := studentType.Field(i)
			field.Tag.Get("json")
			if field.Tag.Get("json") == k+",omitempty" {
				if studentVal.Field(i).CanSet() {
					fieldVal := studentVal.Field(i)
					fieldVal.Set(reflect.ValueOf(v).Convert(studentVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?",
		existingStudent.FirstName, existingStudent.LastName, existingStudent.Email, existingStudent.Class, existingStudent.ID)
	if err != nil {
		log.Println("update student error:", err)

		return model.Student{}, utils.ErrorHandler(err, "Error update data")
	}
	return existingStudent, nil
}

func DeleteOneStudent(id int) error {
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)

		return utils.ErrorHandler(err, "Error update data")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {

		return utils.ErrorHandler(err, "Error update data")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {

		return utils.ErrorHandler(err, "Error update data")
	}

	if rowsAffected == 0 {

		return utils.ErrorHandler(err, "student not found")
	}
	return nil
}

func DeleteStudents(ids []int) ([]int, error) {
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

	stmt, err := tx.Prepare("DELETE FROM students WHERE id = ?")
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
