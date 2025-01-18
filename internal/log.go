package internal

import (
	"context"
	"net/http"

	"github.com/immanuel-254/myauth/internal/models"
)

func LogList(w http.ResponseWriter, r *http.Request) {
	queries := models.New(DB)
	ctx := context.Background()

	auth := ctx.Value("current_user")

	authUser := auth.(models.User)

	logs, err := queries.LogList(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logging(queries, ctx, "log", "list", 0, authUser.ID, w, r)

	SendData(map[string]interface{}{"logs": logs}, w, r)
}
