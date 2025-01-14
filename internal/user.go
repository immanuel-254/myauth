package internal

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/immanuel-254/myauth/internal/models"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	queries := models.New(DB)
	ctx := context.Background()

	// get data
	var data map[string]string
	GetData(data, w, r)

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
	user, err := queries.UserCreate(ctx, models.UserCreateParams{
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

	// send email
	one_time, err := GenerateOneTimeToken(32, uint(user.ID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendEmail(user.Email, "Activate Your Email", fmt.Sprintf("%s/activate/?token=%s", os.Getenv("DOMAIN"), one_time), EmailVerificationTemplate, w, r)

	resp := map[string]interface{}{"message": "signup successful"}
	SendData(resp, w, r)

}

func ActivateEmail(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	token := queryParams.Get("token")

	// verify token
	user_id, err := VerifyToken(token)

	if err != nil {
		http.Error(w, "Invalid auth token", http.StatusBadRequest)
		return
	}

	queries := models.New(DB)
	ctx := context.Background()

	// activate user
	_, err = queries.UserUpdateIsActive(ctx, models.UserUpdateIsActiveParams{
		ID:        int64(user_id),
		Isactive:  sql.NullBool{Bool: true, Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{"message": "email has been verified"}
	SendData(resp, w, r)
}

// Require auth
func ChangeEmailRequest(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// get data
	var data map[string]string
	GetData(data, w, r)

	auth := ctx.Value("current_user")

	authUser := auth.(models.User)

	if data["email"] != authUser.Email {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	// send email
	one_time, err := GenerateOneTimeToken(32, uint(authUser.ID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendEmail(data["email"], "Change Your Email", fmt.Sprintf("%s/change-email/?token=%s", os.Getenv("DOMAIN"), one_time), ChangeEmailVerificationTemplate, w, r)

	SendData(map[string]interface{}{"message": "email sent"}, w, r)

}

func ChangeEmail(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	token := queryParams.Get("token")

	// verify token
	user_id, err := VerifyToken(token)

	if err != nil {
		http.Error(w, "Invalid auth token", http.StatusBadRequest)
		return
	}

	queries := models.New(DB)
	ctx := context.Background()

	auth := ctx.Value("current_user")

	authUser := auth.(models.User)

	if authUser.ID != int64(user_id) {
		http.Error(w, "Forbidden User", http.StatusForbidden)
		return
	}

	// get data
	var data map[string]string
	GetData(data, w, r)

	_, err = queries.UserUpdateEmail(ctx, models.UserUpdateEmailParams{
		ID:        int64(user_id),
		Email:     data["email"],
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendData(map[string]interface{}{"message": "email updated successfully"}, w, r)
}

func ChangePasswordRequest(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	auth := ctx.Value("current_user")

	authUser := auth.(models.User)

	// send email
	one_time, err := GenerateOneTimeToken(32, uint(authUser.ID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendEmail(authUser.Email, "Change Your Password", fmt.Sprintf("%s/change-password/?token=%s", os.Getenv("DOMAIN"), one_time), ChangePasswordVerificationTemplate, w, r)

	SendData(map[string]interface{}{"message": "email sent"}, w, r)

}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	token := queryParams.Get("token")

	// verify token
	user_id, err := VerifyToken(token)

	if err != nil {
		http.Error(w, "Invalid auth token", http.StatusBadRequest)
		return
	}

	queries := models.New(DB)
	ctx := context.Background()

	auth := ctx.Value("current_user")

	authUser := auth.(models.User)

	if authUser.ID != int64(user_id) {
		http.Error(w, "Forbidden User", http.StatusForbidden)
		return
	}

	// get data
	var data map[string]string
	GetData(data, w, r)

	check := CheckPasswordHash(data["old_password"], authUser.Password)

	if !check {
		http.Error(w, "Invalid Password", http.StatusBadRequest)
		return
	}

	if data["new_password"] != data["confirm_password"] {
		http.Error(w, "Invalid Password", http.StatusBadRequest)
		return
	}

	hash, err := HashPassword(data["new_password"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = queries.UserUpdatePassword(ctx, models.UserUpdatePasswordParams{
		ID:        int64(user_id),
		Password:  hash,
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendData(map[string]interface{}{"message": "password updated successfully"}, w, r)
}

func ResetPasswordRequest(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	auth := ctx.Value("current_user")

	authUser := auth.(models.User)

	// send email
	one_time, err := GenerateOneTimeToken(32, uint(authUser.ID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendEmail(authUser.Email, "Reset Your Password", fmt.Sprintf("%s/reset-password/?token=%s", os.Getenv("DOMAIN"), one_time), ResetPasswordVerificationTemplate, w, r)

	SendData(map[string]interface{}{"message": "email sent"}, w, r)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	token := queryParams.Get("token")

	// verify token
	user_id, err := VerifyToken(token)

	if err != nil {
		http.Error(w, "Invalid auth token", http.StatusBadRequest)
		return
	}

	queries := models.New(DB)
	ctx := context.Background()

	auth := ctx.Value("current_user")

	authUser := auth.(models.User)

	if authUser.ID != int64(user_id) {
		http.Error(w, "Forbidden User", http.StatusForbidden)
		return
	}

	// get data
	var data map[string]string
	GetData(data, w, r)

	if data["new_password"] != data["confirm_password"] {
		http.Error(w, "Invalid Password", http.StatusBadRequest)
		return
	}

	hash, err := HashPassword(data["new_password"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = queries.UserUpdatePassword(ctx, models.UserUpdatePasswordParams{
		ID:        int64(user_id),
		Password:  hash,
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendData(map[string]interface{}{"message": "password updated successfully"}, w, r)
}

func DeleteUserRequest(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	auth := ctx.Value("current_user")

	authUser := auth.(models.User)

	// send email
	one_time, err := GenerateOneTimeToken(32, uint(authUser.ID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendEmail(authUser.Email, "Delete User Account", fmt.Sprintf("%s/delete-user/?token=%s", os.Getenv("DOMAIN"), one_time), DeleteUserVerificationTemplate, w, r)

	SendData(map[string]interface{}{"message": "email sent"}, w, r)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	token := queryParams.Get("token")

	// verify token
	user_id, err := VerifyToken(token)

	if err != nil {
		http.Error(w, "Invalid auth token", http.StatusBadRequest)
		return
	}

	queries := models.New(DB)
	ctx := context.Background()

	auth := ctx.Value("current_user")

	authUser := auth.(models.User)

	if authUser.ID != int64(user_id) {
		http.Error(w, "Forbidden User", http.StatusForbidden)
		return
	}

	err = queries.UserDelete(ctx, int64(user_id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendData(map[string]interface{}{"message": "user account deleted"}, w, r)
}

// require admin
func IsActiveChange(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	user_id, err := strconv.ParseInt(queryParams.Get("token"), 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queries := models.New(DB)
	ctx := context.Background()

	// get data
	var data map[string]string
	GetData(data, w, r)

	status, err := strconv.ParseBool(data["active"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = queries.UserUpdateIsActive(ctx, models.UserUpdateIsActiveParams{
		ID:        user_id,
		Isactive:  sql.NullBool{Bool: status, Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendData(map[string]interface{}{"message": "user active status updated successfully"}, w, r)
}

func IsStaffChange(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	user_id, err := strconv.ParseInt(queryParams.Get("token"), 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queries := models.New(DB)
	ctx := context.Background()

	// get data
	var data map[string]string
	GetData(data, w, r)

	status, err := strconv.ParseBool(data["staff"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = queries.UserUpdateIsStaff(ctx, models.UserUpdateIsStaffParams{
		ID:        user_id,
		Isstaff:   sql.NullBool{Bool: status, Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendData(map[string]interface{}{"message": "user staff status updated successfully"}, w, r)
}
