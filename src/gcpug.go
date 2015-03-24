package gcpug

import (
	"net/http"
)

func init() {
	http.HandleFunc("/gaego", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
}
