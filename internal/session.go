package internal

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/immanuel-254/myauth/internal/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	queries := models.New(DB)
	ctx := context.Background()

	// get data
	var data map[string]string
	GetData(data, w, r)

	// get user
	user, err := queries.UserLoginRead(ctx, data["email"])

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err)
		return
	}

	check := CheckPasswordHash(data["password"], user.Password)

	if !check {
		http.Error(w, "Invalid Credentials", http.StatusBadRequest)
		return
	}

	// create key
	key := string(GenerateAESKey())

	// create session

	session, err := queries.SessionCreate(ctx, models.SessionCreateParams{
		Key:       key,
		UserID:    user.ID,
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = queries.LogCreate(ctx, models.LogCreateParams{
		DbTable:   "session",
		Action:    "create",
		ObjectID:  session.ID,
		UserID:    session.UserID,
		CreatedAt: sql.NullTime{Time: time.Now()},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{"auth": key}
	SendData(resp, w, r)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	queries := models.New(DB)
	ctx := context.Background()

	// delete session
	err := queries.SessionDelete(ctx, w.Header().Get("auth"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	auth := ctx.Value("current_user")

	authUser := auth.(models.User)

	err = queries.LogCreate(ctx, models.LogCreateParams{
		DbTable:   "session",
		Action:    "delete",
		ObjectID:  0,
		UserID:    authUser.ID,
		CreatedAt: sql.NullTime{Time: time.Now()},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{"message": "user deleted"}
	SendData(resp, w, r)
}

func SessionList(w http.ResponseWriter, r *http.Request) {
	queries := models.New(DB)
	ctx := context.Background()

	auth := ctx.Value("current_user")

	authUser := auth.(models.User)

	sessions, err := queries.SessionList(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logging(queries, ctx, "session", "list", 0, authUser.ID, w, r)

	SendData(map[string]interface{}{"sessions": sessions}, w, r)
}
