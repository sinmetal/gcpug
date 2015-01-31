package gcpug

import (
	"fmt"
	"net/http"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func init() {
	route(goji.DefaultMux)
	goji.Serve()
}

func route(m *web.Mux) {
	m.Get("/hello/:name", hello)
}

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}
