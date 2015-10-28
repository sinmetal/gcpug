package gcpug

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/zenazn/goji/web"
	"golang.org/x/net/context"

	"code.google.com/p/go-uuid/uuid"
	"github.com/mjibson/goon"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type PugEvent struct {
	Id             string    `datastore:"-" goon:"id" json:"id"`       // UUID
	OrganizationId string    `json:"organizationId"`                   // 支部Id
	Title          string    `json:"title" datastore:",noindex"`       // イベントタイトル
	Description    string    `json:"description" datastore:",noindex"` // イベント説明
	Url            string    `json:"url" datastore:",noindex"`         // イベント募集URL
	StartAt        time.Time `json:"startAt"`                          // 開催日時
	CreatedAt      time.Time `json:"createdAt"`                        // 作成日時
	UpdatedAt      time.Time `json:"updatedAt"`                        // 更新日時
}

type PugEventApi struct {
}

func SetUpPugEvent(m *web.Mux) {
	api := PugEventApi{}

	m.Get("/api/1/event", api.List)
	m.Post("/api/1/event", api.Post)
	m.Put("/api/1/event", api.Put)
}

func (a *PugEventApi) Post(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var e PugEvent
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		log.Infof(ctx, "request decode error : %s", err.Error())
		er := ErrorResponse{
			http.StatusBadRequest,
			[]string{"invalid request"},
		}
		er.Write(w)
		return
	}
	defer r.Body.Close()

	g := goon.NewGoon(r)
	e.Id = uuid.New()
	err = e.Create(g)
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
	json.NewEncoder(w).Encode(e)
}

func (a *PugEventApi) Put(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var e PugEvent
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		log.Infof(ctx, "request decode error : %s", err.Error())
		er := ErrorResponse{
			http.StatusBadRequest,
			[]string{"invalid request"},
		}
		er.Write(w)
		return
	}
	defer r.Body.Close()

	err = e.Validate4Put()
	if err != nil {
		er := ErrorResponse{
			http.StatusBadRequest,
			[]string{err.Error()},
		}
		er.Write(w)
		return
	}

	err = e.Update(ctx)
	if err != nil {
		er := ErrorResponse{
			http.StatusNotFound,
			[]string{err.Error()},
		}
		er.Write(w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(e)
}

func (a *PugEventApi) List(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	g := goon.FromContext(ctx)

	q := datastore.NewQuery(goon.DefaultKindName(&PugEvent{}))
	q = q.Order("-StartAt")

	eok := r.FormValue("organizationKey")
	if eok != "" {
		k, err := datastore.DecodeKey(eok)
		if err != nil {
			log.Infof(ctx, "invalid organizationKey : %s", eok)
			er := ErrorResponse{
				http.StatusBadRequest,
				[]string{"invalid organizationKey"},
			}
			er.Write(w)
			return
		}
		q = q.Filter("OrganizationKey = ", k)
	}

	limit := r.FormValue("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			log.Infof(ctx, "invalid limit : %s", limit)
			er := ErrorResponse{
				http.StatusBadRequest,
				[]string{"invalid limit"},
			}
			er.Write(w)
			return
		}
		q = q.Limit(l)
	}

	pes := make([]*PugEvent, 0)
	_, err := g.GetAll(q, &pes)
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
	json.NewEncoder(w).Encode(pes)
}

func (pe *PugEvent) Validate() error {
	if pe.Title == "" {
		return errors.New("title is required.")
	}
	return nil
}

func (pe *PugEvent) Validate4Put() error {
	if pe.Id == "" {
		return errors.New("id is required.")
	}
	return pe.Validate()
}

func (pe *PugEvent) Get(c context.Context) error {
	g := goon.FromContext(c)
	return g.Get(pe)
}

func (pe *PugEvent) Update(c context.Context) error {
	g := goon.FromContext(c)
	err := g.RunInTransaction(func(g *goon.Goon) error {
		stored := &PugEvent{
			Id: pe.Id,
		}
		err := g.Get(stored)
		if err != nil {
			return err
		}

		pe.CreatedAt = stored.CreatedAt

		_, err = g.Put(pe)
		if err != nil {
			return err
		}

		return nil
	}, nil)
	if err != nil {
		log.Warningf(c, "%v", pe)
		return err
	}

	j, err := json.Marshal(pe)
	if err != nil {
		log.Warningf(c, "json marshal error, %s", pe.Id)
	}
	log.Infof(c, "{\"__DS__KIND__PUGEVENT__\":%s}", j)
	return err
}

func (pe *PugEvent) Load(ps []datastore.Property) error {
	if err := datastore.LoadStruct(pe, ps); err != nil {
		return err
	}
	return nil
}

func (pe *PugEvent) Save() ([]datastore.Property, error) {
	now := time.Now()
	pe.UpdatedAt = now

	if pe.CreatedAt.IsZero() {
		pe.CreatedAt = now
	}

	return datastore.SaveStruct(pe)
}

func (pe *PugEvent) Create(g *goon.Goon) error {
	return g.RunInTransaction(func(g *goon.Goon) error {
		stored := &PugEvent{
			Id: pe.Id,
		}
		err := g.Get(stored)
		if err == nil {
			return ConflictKey
		}
		if err != datastore.ErrNoSuchEntity {
			return err
		}

		_, err = g.Put(pe)
		if err != nil {
			return err
		}

		return nil
	}, nil)
}
