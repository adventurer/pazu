package middleware

import (
	"net/http"
)

func Third(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Write([]byte("third"))
	next(w, r)
	return

}
