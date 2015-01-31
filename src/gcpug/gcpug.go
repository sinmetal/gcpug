package gcpug

import (
	"encoding/json"
	"fmt"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"net/http"
	"time"
)

type Organization struct {
	Id        string    `json:id`
	Name      string    `json:name`
	Url       string    `json:url`
	CreatedAt time.Time `json createdAt`
}

func init() {
	route(goji.DefaultMux)
	goji.Serve()
}

func route(m *web.Mux) {
	m.Get("/hello/:name", hello)
	m.Get("/organization/:id", doGetOrganization)
	m.Get("/organization", doGetOrganizationList)
}

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

func doGetOrganization(c web.C, w http.ResponseWriter, r *http.Request) {
	o := Organization{
		"sampleid",
		"Sinmetal支部",
		"http://sinmetal.org",
		time.Now(),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(o)
}

func doGetOrganizationList(c web.C, w http.ResponseWriter, r *http.Request) {

	o := []Organization{
		Organization{
			"sampleid1",
			"Sinmetal支部1",
			"http://sinmetal1.org",
			time.Now(),
		},
		Organization{
			"sampleid2",
			"Sinmetal支部2",
			"http://sinmetal2.org",
			time.Now(),
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(o)
}
