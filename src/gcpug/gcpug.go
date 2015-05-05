package gcpug

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/zenazn/goji"
)

var (
	ConflictKey = errors.New("datastore: conflict key")
)

type ErrorResponse struct {
	Status   int
	Messages []string
}

func init() {
	m := goji.DefaultMux

	SetUpOrganization(m)
	SetUpPugEvent(m)

	goji.Serve()
}

func (er *ErrorResponse) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(er.Status)
	json.NewEncoder(w).Encode(er.Messages)
}
