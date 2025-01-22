package internal

import (
	"context"
	"net/http"

	"github.com/immanuel-254/myauth/frontend/src"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	component := src.Base(src.DashBoard())

	component.Render(context.Background(), w)
}
