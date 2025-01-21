package internal

import (
	"context"
	"net/http"

	"github.com/immanuel-254/myauth/frontend/src"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	component := src.Base(src.HelloWorld())

	component.Render(context.Background(), w)
}
