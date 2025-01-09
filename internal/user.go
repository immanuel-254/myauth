package internal

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/immanuel-254/myauth/internal/models"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	queries := models.New(DB)
	ctx := context.Background()

	// get data
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err)
		return
	}

	// check if password match
	if data["password"] != data["confirm-password"] {
		http.Error(w, "invalid password or email", http.StatusBadRequest)
		return
	}

	// hash password
	hash, err := HashPassword(data["password"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create user
	_, err = queries.UserCreate(ctx, models.UserCreateParams{
		Email:     data["email"],
		Password:  hash,
		Isactive:  sql.NullBool{Bool: false, Valid: true},
		Isstaff:   sql.NullBool{Bool: false, Valid: true},
		Isadmin:   sql.NullBool{Bool: false, Valid: true},
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"message": "signup successful"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
func ActivateEmail(w http.ResponseWriter, r *http.Request) {}

// Require auth
func ChangeEmailRequest(w http.ResponseWriter, r *http.Request)    {}
func ChangeEmail(w http.ResponseWriter, r *http.Request)           {}
func ChangePasswordRequest(w http.ResponseWriter, r *http.Request) {}
func ChangePassword(w http.ResponseWriter, r *http.Request)        {}
func ResetPasswordRequest(w http.ResponseWriter, r *http.Request)  {}
func ResetPassword(w http.ResponseWriter, r *http.Request)         {}
func DeleteUser(w http.ResponseWriter, r *http.Request)            {}

// require admin
func CreateStaff(w http.ResponseWriter, r *http.Request)    {}
func IsActiveChange(w http.ResponseWriter, r *http.Request) {}
func IsStaffChange(w http.ResponseWriter, r *http.Request)  {}
