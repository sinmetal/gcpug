package gcpug

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"appengine"
	"appengine/datastore"
	"github.com/mjibson/goon"
)

var (
	ConflictKey = errors.New("datastore: conflict key")
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
	Status   int
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
	m.Post("/api/1/organization", api.Post)
}

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

func (a *OrganizationApi) Get(c web.C, w http.ResponseWriter, r *http.Request) {
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
		ac.Errorf(fmt.Sprintf("datastore get error : ", err.Error()))
		er.Write(w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(o)
}

func (a *OrganizationApi) list(c web.C, w http.ResponseWriter, r *http.Request) {
	ac := appengine.NewContext(r)
	g := goon.FromContext(ac)

	q := datastore.NewQuery(goon.DefaultKindName(&Organization{}))
	q.Order("CreatedAt")

	os := make([]*Organization, 0)
	_, err := g.GetAll(q, &os)
	if err != nil {
		ac.Errorf(err.Error())
		er := ErrorResponse{
			http.StatusInternalServerError,
			[]string{"datastore query error"},
		}
		er.Write(w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(os)
}

func (a *OrganizationApi) Post(c web.C, w http.ResponseWriter, r *http.Request) {
	ac := appengine.NewContext(r)

	var o Organization
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		ac.Infof("rquest body, %v", r.Body)
		er := ErrorResponse{
			http.StatusBadRequest,
			[]string{"invalid request"},
		}
		er.Write(w)
		return
	}
	defer r.Body.Close()

	err = o.Create(ac)
	if err == ConflictKey {
		er := ErrorResponse{
			http.StatusBadRequest,
			[]string{"conflict Id"},
		}
		er.Write(w)
		return
	} else if err != nil {
		er := ErrorResponse{
			http.StatusInternalServerError,
			[]string{err.Error()},
		}
		ac.Errorf(fmt.Sprintf("datastore put error : ", err.Error()))
		er.Write(w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(o)
}

func (o *Organization) Get(c appengine.Context) error {
	g := goon.FromContext(c)
	return g.Get(o)
}

func (o *Organization) Create(c appengine.Context) error {
	g := goon.FromContext(c)
	return g.RunInTransaction(func(g *goon.Goon) error {
		stored := &Organization{
			Id: o.Id,
		}
		err := g.Get(stored)
		if err != datastore.ErrNoSuchEntity {
			return ConflictKey
		}

		_, err = g.Put(o)
		if err != nil {
			return err
		}

		return nil
	}, nil)
}

func (o *Organization) Load(c <-chan datastore.Property) error {
	if err := datastore.LoadStruct(o, c); err != nil {
		return err
	}

	return nil
}

func (o *Organization) Save(c chan<- datastore.Property) error {
	now := time.Now()
	o.UpdatedAt = now

	if o.CreatedAt.IsZero() {
		o.CreatedAt = now
	}

	if err := datastore.SaveStruct(o, c); err != nil {
		return err
	}
	return nil
}

func (er *ErrorResponse) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(er.Status)
	json.NewEncoder(w).Encode(er.Messages)
}
