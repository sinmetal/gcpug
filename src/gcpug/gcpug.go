package gcpug

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"appengine"
	"appengine/datastore"
	"github.com/mjibson/goon"
)

// Organization
//
// 支部
type Organization struct {
	Id        string    `datastore:"-" goon:"id" json:id` // 明示的に入れるID
	Name      string    `json:name datastore:",noindex"`  // 支部名
	Url       string    `json:url datastore:",noindex"`   // 支部WebSiteURL
	CreatedAt time.Time `json createdAt`                  // 作成日時
	UpdatedAt time.Time `json updatedAt`                  // 更新日時
}

type OrganizationApi struct {
}

type ErrorResponse struct {
	Status int
	Messages []string
}

func init() {
	route(goji.DefaultMux)
	goji.Serve()
}

func route(m *web.Mux) {
	api := OrganizationApi{}

	m.Get("/hello/:name", hello)
	m.Get("/api/1/organization/:id", api.Get)
	m.Get("/api/1/organization", api.list)
}

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

func (a *OrganizationApi) Get (c web.C, w http.ResponseWriter, r *http.Request) {
	ac := appengine.NewContext(r)

	id := c.URLParams["id"]
	if id == "" {
		er := ErrorResponse{
			http.StatusBadRequest,
			[]string{"id is required"},
		}
		er.Write(w)
		return
	}

	o := &Organization{Id: id}
	err := o.Get(ac)
	if err == datastore.ErrNoSuchEntity {
		er := ErrorResponse{
			http.StatusNotFound,
			[]string{fmt.Sprintf("%s is not found.", id)},
		}
		er.Write(w)
		return
	} else if err != nil {
		er := ErrorResponse{
			http.StatusInternalServerError,
			[]string{err.Error()},
		}
		ac.Errorf(fmt.Sprintf("datastore get error : ", err.Error() ))
		er.Write(w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(o)
}

func (a *OrganizationApi) list(c web.C, w http.ResponseWriter, r *http.Request) {

	o := []Organization{
		Organization{
			"sampleid1",
			"Sinmetal支部1",
			"http://sinmetal1.org",
			time.Now(),
			time.Now(),
		},
		Organization{
			"sampleid2",
			"Sinmetal支部2",
			"http://sinmetal2.org",
			time.Now(),
			time.Now(),
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(o)
}

func (o *Organization) Get(c appengine.Context) error {
	g := goon.FromContext(c)
	return g.Get(o)
}

func (er *ErrorResponse) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(er.Status)
	json.NewEncoder(w).Encode(er.Messages)
}
