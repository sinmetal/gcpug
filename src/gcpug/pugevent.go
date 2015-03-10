package gcpug

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"appengine"
	"appengine/datastore"
	"code.google.com/p/go-uuid/uuid"
	"github.com/mjibson/goon"
)

type PugEvent struct {
	Id              string         `datastore:"-" goon:"id" json:"id"` // UUID
	OrganizationKey *datastore.Key `json:"organizationKey"`            // 支部KindKey
	Title           string         `json:"title" datastore:",noindex"` // イベントタイトル
	Url             string         `json:"url" datastore:",noindex"`   // イベント募集URL
	StartAt         time.Time      `json:"startAt"`                    // 開催日時
	CreatedAt       time.Time      `json:"createdAt"`                  // 作成日時
	UpdatedAt       time.Time      `json:"updatedAt"`                  // 更新日時
}

type PugEventApi struct {
}

func SetUpPugEvent(m *web.Mux) {
	api := PugEventApi{}

	m.Post("/api/1/event", api.Post)
}

func (a *PugEventApi) Post(c web.C, w http.ResponseWriter, r *http.Request) {
	ac := appengine.NewContext(r)

	var e PugEvent
	err := json.NewDecoder(r.Body).Decode(&e)
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

	e.Id = uuid.New()
	err = e.Create(ac)
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
	json.NewEncoder(w).Encode(e)
}

func (pe *PugEvent) Validate() error {
	if pe.Title == "" {
		return errors.New("title is required.")
	}
	return nil
}

func (pe *PugEvent) Load(c <-chan datastore.Property) error {
	if err := datastore.LoadStruct(pe, c); err != nil {
		return err
	}

	return nil
}

func (pe *PugEvent) Save(c chan<- datastore.Property) error {
	now := time.Now()
	pe.UpdatedAt = now

	if pe.CreatedAt.IsZero() {
		pe.CreatedAt = now
	}

	if err := datastore.SaveStruct(pe, c); err != nil {
		return err
	}
	return nil
}

func (pe *PugEvent) Create(c appengine.Context) error {
	g := goon.FromContext(c)
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
