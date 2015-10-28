package gcpug

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/zenazn/goji/web"

	"github.com/mjibson/goon"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// Organization
//
// 支部
type Organization struct {
	Id        string    `datastore:"-" goon:"id" json:"id"`   // 明示的に入れるID
	Name      string    `json:"name" datastore:",noindex"`    // 支部名
	Url       string    `json:"url" datastore:",noindex"`     // 支部WebSiteURL
	LogoUrl   string    `json:"logoUrl" datastore:",noindex"` // LogoURL
	Order     int       `json:"order"`                        // 並び順
	CreatedAt time.Time `json:"createdAt"`                    // 作成日時
	UpdatedAt time.Time `json:"updatedAt"`                    // 更新日時
}

type OrganizationApi struct {
}

func SetUpOrganization(m *web.Mux) {
	api := OrganizationApi{}

	m.Get("/api/1/organization/:id", api.Get)
	m.Get("/api/1/organization", api.List)
	m.Post("/api/1/organization", api.Post)
	m.Put("/api/1/organization", api.Put)
}

func (a *OrganizationApi) Get(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

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
	err := o.Get(ctx)
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
		log.Errorf(ctx, fmt.Sprintf("datastore get error : ", err.Error()))
		er.Write(w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(o)
}

func (a *OrganizationApi) List(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	g := goon.FromContext(ctx)

	q := datastore.NewQuery(goon.DefaultKindName(&Organization{})).
		Order("Order")

	os := make([]*Organization, 0)
	_, err := g.GetAll(q, &os)
	if err != nil {
		log.Errorf(ctx, err.Error())
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
	ctx := appengine.NewContext(r)

	var o Organization
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		log.Infof(ctx, "rquest body, %v", r.Body)
		er := ErrorResponse{
			http.StatusBadRequest,
			[]string{"invalid request"},
		}
		er.Write(w)
		return
	}
	defer r.Body.Close()

	err = o.Create(ctx)
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
		log.Errorf(ctx, fmt.Sprintf("datastore put error : ", err.Error()))
		er.Write(w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(o)
}

func (a *OrganizationApi) Put(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var o Organization
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		log.Infof(ctx, "rquest body, %v", r.Body)
		er := ErrorResponse{
			http.StatusBadRequest,
			[]string{"invalid request"},
		}
		er.Write(w)
		return
	}
	defer r.Body.Close()

	err = o.Update(ctx)
	if err != nil {
		er := ErrorResponse{
			http.StatusInternalServerError,
			[]string{err.Error()},
		}
		log.Errorf(ctx, fmt.Sprintf("datastore put error : ", err.Error()))
		er.Write(w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(o)
}

func (o *Organization) Get(c context.Context) error {
	g := goon.FromContext(c)
	return g.Get(o)
}

func (o *Organization) Create(c context.Context) error {
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

func (o *Organization) Update(c context.Context) error {
	g := goon.FromContext(c)
	return g.RunInTransaction(func(g *goon.Goon) error {
		stored := &Organization{
			Id: o.Id,
		}
		err := g.Get(stored)
		if err != nil {
			return err
		}

		o.CreatedAt = stored.CreatedAt
		o.UpdatedAt = stored.UpdatedAt

		log.Infof(c, "new name %s", o.Name)
		_, err = g.Put(o)
		if err != nil {
			return err
		}

		return nil
	}, nil)
}

func (o *Organization) Load(ps []datastore.Property) error {
	if err := datastore.LoadStruct(o, ps); err != nil {
		return err
	}
	return nil
}

func (o *Organization) Save() ([]datastore.Property, error) {
	now := time.Now()
	o.UpdatedAt = now

	if o.CreatedAt.IsZero() {
		o.CreatedAt = now
	}

	return datastore.SaveStruct(o)
}
