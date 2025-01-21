package internal

import "net/http"

func StyleCss(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/styles.css")
}
func AlpineJs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/alpine.js")
}
