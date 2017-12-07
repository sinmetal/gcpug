package gcpug

import (
	"net/http"
	"fmt"

	"github.com/sinmetal/ds2bq"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type DSBUAPI struct {}

func init() {
	api := DSBUAPI{}
	http.HandleFunc("/cron/dsbu", api.Get)
}

func (api *DSBUAPI) Get(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	s := ds2bq.NewDatastoreExportService()
	op, err := s.Export(ctx, fmt.Sprintf("gs://%s-backup", appengine.AppID(ctx)), &ds2bq.EntityFilter{
		Kinds:[]string{
			"Organization",
			"PugEvent",
			"Stackoverflow",
		},
		NamespaceIds: []string{},
	})
	if err != nil {
		log.Errorf(ctx, "%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(op.HTTPStatusCode)
	rj, err := op.Response.MarshalJSON()
	if err != nil {
		log.Errorf(ctx, "%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof(ctx, "%s", string(rj))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(rj)
}


