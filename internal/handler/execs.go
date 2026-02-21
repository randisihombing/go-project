package handler

import (
	"encoding/json"
	"fmt"
	"gocourse/internal/model"
	"gocourse/internal/repository/sqlconnect"
	"gocourse/pkg/utils"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req model.Exec

	//Data Validation
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and Password are required", http.StatusBadRequest)
		return
	}

	//Search for User Already Exist
	user, err := sqlconnect.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "invalid username or password", http.StatusBadRequest)
		return
	}

	//is user active
	if user.InactiveStatus {
		http.Error(w, "account is inactive", http.StatusForbidden)
		return
	}

	//Verify password
	err = utils.VerifyPassword(req.Password, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	//Generate JWT Token
	tokenString, err := utils.SignToken(user.ID, req.Username, user.Role)
	if err != nil {
		http.Error(w, "could not create token", http.StatusInternalServerError)
		return
	}

	//Send token as a response or as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "test",
		Value:    "testing",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})

	//response body
	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}
	json.NewEncoder(w).Encode(response)

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Log out successfully"}`))
}

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid exec ID", http.StatusBadRequest)
		return
	}

	var req model.UpdatePasswordRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	if req.CurrentPassword == "" || req.NewPassword == "" {
		http.Error(w, "Please enter password", http.StatusBadRequest)
		return
	}

	_, err = sqlconnect.UpdatePasswordInDb(userId, req.CurrentPassword, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// //Send token as a response or as a cookie
	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "Bearer",
	// 	Value:    token,
	// 	Path:     "/",
	// 	HttpOnly: true,
	// 	Secure:   true,
	// 	Expires:  time.Now().Add(24 * time.Hour),
	// 	SameSite: http.SameSiteStrictMode,
	// })

	//response body
	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Password has been updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	err = sqlconnect.ForgotPasswordDbHandler(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//respond with success message
	fmt.Fprintf(w, "Password reset link send to %s", req.Email)

}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("resetcode")

	type request struct {
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	var req request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid value in request", http.StatusBadRequest)
		return
	}

	// TODO: data validation for blank values
	if req.NewPassword != req.ConfirmPassword {
		http.Error(w, "Password should match", http.StatusBadRequest)
		return
	}

	err = sqlconnect.ResetPasswordDbHandler(token, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, "Password reset successfully")

}
