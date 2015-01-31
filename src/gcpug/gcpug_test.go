package gcpug

import (
	"testing"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"github.com/zenazn/goji/web"
)

func TestHello(t *testing.T) {
	m := web.New()
	route(m)
	ts := httptest.NewServer(m)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/hello/sinmetal")
	if err != nil {
		t.Error("unexpected")
	}
	if res.StatusCode != http.StatusOK {
		t.Error("invalid status code")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if (string(body) != "Hello, sinmetal!") {
		t.Error("invalid body : ", string(body))
	}
}
