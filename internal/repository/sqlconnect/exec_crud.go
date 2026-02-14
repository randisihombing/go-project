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
)

func GetExecsDbHandler(execs []model.Exec, r *http.Request) ([]model.Exec, error) {
	db, err := ConnectDb()
	if err != nil {

		return nil, utils.ErrorHandler(err, "Error retrieving data")
	}

	defer db.Close()

	query := "SELECT id, first_name, last_name, email, username, user_created_at, inactive_status, role FROM execs WHERE 1=1"

	var args []interface{}

	query, args = utils.AddFilters(r, query, args)

	utils.AddSorting(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)

		return nil, utils.ErrorHandler(err, "Error retrieving data")
	}
	defer rows.Close()

	for rows.Next() {
		var exec model.Exec
		err := rows.Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.UserCreatedAt, &exec.InactiveStatus, &exec.Role)
		if err != nil {

			return nil, utils.ErrorHandler(err, "Error retrieving data")
		}
		execs = append(execs, exec)
	}
	return execs, nil
}

func GetExecByID(id int) (model.Exec, error) {
	db, err := ConnectDb()
	if err != nil {

		return model.Exec{}, utils.ErrorHandler(err, "Error retrieving data")
	}

	defer db.Close()

	var exec model.Exec
	err = db.QueryRow("SELECT id, first_name, last_name, email, username, inactive_status, role FROM execs WHERE id = ?", id).Scan(
		&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.InactiveStatus, &exec.Role)
	if err == sql.ErrNoRows {

		return model.Exec{}, utils.ErrorHandler(err, "Error retrieving data")
	} else if err != nil {

		return model.Exec{}, utils.ErrorHandler(err, "Error retrieving data")
	}
	return exec, nil
}

func AddExecsDbHandler(newExecs []model.Exec) ([]model.Exec, error) {
	db, err := ConnectDb()
	if err != nil {

		return nil, utils.ErrorHandler(err, "Error adding data")
	}

	defer db.Close()

	stmt, err := db.Prepare(utils.GenerateInsertQuery("execs", model.Exec{}))
	if err != nil {

		return nil, utils.ErrorHandler(err, "Error adding data")
	}
	defer stmt.Close()

	addedExecs := make([]model.Exec, len(newExecs))
	for i, newExec := range newExecs {
		values := utils.GetStructValues(newExec)
		res, err := stmt.Exec(values...)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error adding data")
		}
		lastID, err := res.LastInsertId()
		if err != nil {

			return nil, utils.ErrorHandler(err, "Error adding data")
		}
		newExec.ID = int(lastID)

		addedExecs[i] = newExec
		// nextID++
	}
	return addedExecs, nil
}

func PatchExecs(updates []map[string]interface{}) error {
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

		var execFromDb model.Exec
		err = db.QueryRow("SELECT id, first_name, last_name, email, username FROM execs WHERE id = ?", id).Scan(
			&execFromDb.ID, &execFromDb.FirstName, &execFromDb.LastName, &execFromDb.Email, &execFromDb.Username)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {

				return utils.ErrorHandler(err, "Exec not found")
			}

			return utils.ErrorHandler(err, "Error update data")
		}

		//Apply update using reflection
		execVal := reflect.ValueOf(&execFromDb).Elem()
		execType := execVal.Type()

		for k, v := range update {
			if k == "id" {
				//skip updating id field
				continue
			}
			for i := 0; i < execVal.NumField(); i++ {
				field := execType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := execVal.Field(i)
					if execVal.Field(i).CanSet() {
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

		_, err = tx.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ? WHERE id = ?",
			execFromDb.FirstName, execFromDb.LastName, execFromDb.Email, execFromDb.ID, execFromDb.Username)
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

func PatchOneExec(id int, updates map[string]interface{}) (model.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)

		return model.Exec{}, utils.ErrorHandler(err, "Error update data")
	}
	defer db.Close()

	var existingExec model.Exec
	err = db.QueryRow("SELECT id, first_name, last_name, email, username FROM execs WHERE id = ?", id).Scan(
		&existingExec.ID, &existingExec.FirstName, &existingExec.LastName, &existingExec.Email, &existingExec.Username)

	if err != nil {
		if err == sql.ErrNoRows {

			return model.Exec{}, utils.ErrorHandler(err, "Exect not found")
		}

		return model.Exec{

			//Reflect package instead switch case
		}, utils.ErrorHandler(err, "Error update data")
	}

	execVal := reflect.ValueOf(&existingExec).Elem()
	execType := execVal.Type()

	for k, v := range updates {
		for i := 0; i < execVal.NumField(); i++ {
			field := execType.Field(i)
			field.Tag.Get("json")
			if field.Tag.Get("json") == k+",omitempty" {
				if execVal.Field(i).CanSet() {
					fieldVal := execVal.Field(i)
					fieldVal.Set(reflect.ValueOf(v).Convert(execVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ? WHERE id = ?",
		existingExec.FirstName, existingExec.LastName, existingExec.Email, existingExec.Username, existingExec.ID)
	if err != nil {
		log.Println("update exec error:", err)

		return model.Exec{}, utils.ErrorHandler(err, "Error update data")
	}
	return existingExec, nil
}

func DeleteOneExec(id int) error {
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)

		return utils.ErrorHandler(err, "Error update data")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM execs WHERE id = ?", id)
	if err != nil {

		return utils.ErrorHandler(err, "Error update data")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {

		return utils.ErrorHandler(err, "Error update data")
	}

	if rowsAffected == 0 {

		return utils.ErrorHandler(err, "exec not found")
	}
	return nil
}
