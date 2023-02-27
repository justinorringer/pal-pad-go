package endpoints

import (
	"net/http"
)

func Lubdub(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("lubdub"))

	w.WriteHeader(http.StatusOK)
}
