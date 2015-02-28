package gcpug

import (
	"encoding/json"
	"github.com/zenazn/goji/web"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mjibson/goon"

	"github.com/sinmetal/gaego_unittest_util/aetestutil"
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
	if string(body) != "Hello, sinmetal!" {
		t.Error("invalid body : ", string(body))
	}
}

func TestDoGetOrganization(t *testing.T) {
	inst, c, err := aetestutil.SpinUp()
	if err != nil {
		t.Fatal(err)
	}
	defer aetestutil.SpinDown()

	g := goon.FromContext(c)

	o := &Organization{
		"sampleId",
		"Sinmetal支部",
		"http://sinmetal.org",
		time.Now(),
		time.Now(),
	}
	_, err = g.Put(o)
	if err != nil {
		t.Fatal("test organization put error : %s", err.Error())
	}

	m := web.New()
	route(m)
	ts := httptest.NewServer(m)
	defer ts.Close()

	r, err := inst.NewRequest("GET", ts.URL + "/api/1/organization/" + o.Id, nil)
	if err != nil {
		t.Error(err.Error())
	}
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Error("unexpected status code : %d, %s", w.Code, w.Body)
	}

	var ro Organization
	json.NewDecoder(w.Body).Decode(&ro)
	if ro.Id != o.Id {
		t.Error("invalid organization.id : ", ro.Id)
	}
	if ro.Name != o.Name {
		t.Error("invalid organization.name : ", ro.Name)
	}
	if ro.Url != o.Url {
		t.Error("invalid organization.url : ", o.Url)
	}
	zeroTime := time.Time{}
	if ro.CreatedAt == zeroTime {
		t.Error("invalid organization.createdAt : ", ro.CreatedAt)
	}
	if ro.UpdatedAt == zeroTime {
		t.Error("invalid organization.UpdatedAt : ", ro.UpdatedAt)
	}
}

func TestDoGetOrganizationList(t *testing.T) {
	m := web.New()
	route(m)
	ts := httptest.NewServer(m)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/api/1/organization")
	if err != nil {
		t.Error("unexpected")
	}
	if res.StatusCode != http.StatusOK {
		t.Error("invalid status code")
	}

	zeroTime := time.Time{}

	defer res.Body.Close()
	var o []Organization
	json.NewDecoder(res.Body).Decode(&o)
	if len(o) != 2 {
		t.Error("invalid organization len : ", len(o))
	}
	if o[0].Id != "sampleid1" {
		t.Error("invalid organization.id : ", o[0].Id)
	}
	if o[0].Name != "Sinmetal支部1" {
		t.Error("invalid organization.name : ", o[0].Name)
	}
	if o[0].Url != "http://sinmetal1.org" {
		t.Error("invalid organization.url : ", o[0].Url)
	}

	if o[0].CreatedAt == zeroTime {
		t.Error("invalid organization.createdAt : ", o[0].CreatedAt)
	}

	if o[1].Id != "sampleid2" {
		t.Error("invalid organization.id : ", o[1].Id)
	}
	if o[1].Name != "Sinmetal支部2" {
		t.Error("invalid organization.name : ", o[1].Name)
	}
	if o[1].Url != "http://sinmetal2.org" {
		t.Error("invalid organization.url : ", o[1].Url)
	}

	if o[1].CreatedAt == zeroTime {
		t.Error("invalid organization.createdAt : ", o[1].CreatedAt)
	}
}
