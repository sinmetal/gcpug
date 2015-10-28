package gcpug

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"github.com/mjibson/goon"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type PugConfig struct {
	Id           string    `datastore:"-" goon:"id" json:"id"`        // pug-config-id 固定
	ClientId     string    `json:"clientId" datastore:",noindex"`     // GCP Client Id
	ClientSecret string    `json:"clientSecret" datastore:",noindex"` // GCP Client Secret
	SlackPostUrl string    `json:"slackPostUrl" datastore:",noindex"` // Slackにぶっこむ用URL
	CreatedAt    time.Time `json:"createdAt"`                         // 作成日時
	UpdatedAt    time.Time `json:"updatedAt"`                         // 更新日時
}

const (
	pugConfigId = "pug-config-id"
)

type PugConfigApi struct {
}

type PugConfigService struct {
}

func SetUpPugConfig(m *web.Mux) {
	api := PugConfigApi{}

	m.Post("/admin/api/1/config", api.Put)
}

func (a *PugConfigApi) Put(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	g := goon.NewGoon(r)

	var pc PugConfig
	err := json.NewDecoder(r.Body).Decode(&pc)
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

	pc.Id = pugConfigId
	_, err = g.Put(&pc)
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
	json.NewEncoder(w).Encode(pc)
}

func (s *PugConfigService) Get(g *goon.Goon) (PugConfig, error) {
	pc := PugConfig{
		Id: pugConfigId,
	}

	err := g.Get(&pc)
	return pc, err
}

func (pc *PugConfig) Load(ps []datastore.Property) error {
	if err := datastore.LoadStruct(pc, ps); err != nil {
		return err
	}
	return nil
}

func (pc *PugConfig) Save() ([]datastore.Property, error) {
	now := time.Now()
	pc.UpdatedAt = now

	if pc.CreatedAt.IsZero() {
		pc.CreatedAt = now
	}

	return datastore.SaveStruct(pc)
}
