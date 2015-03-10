package gcpug

import (
	"bytes"
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

	r, err := inst.NewRequest("GET", ts.URL+"/api/1/organization/"+o.Id, nil)
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
	inst, c, err := aetestutil.SpinUp()
	if err != nil {
		t.Fatal(err)
	}
	defer aetestutil.SpinDown()

	g := goon.FromContext(c)

	o1 := &Organization{
		"sampleId1",
		"Sinmetal支部1",
		"http://sinmetal1.org",
		time.Now(),
		time.Now(),
	}
	_, err = g.Put(o1)
	if err != nil {
		t.Fatal("test organization put error : %s", err.Error())
	}

	o2 := &Organization{
		"sampleId2",
		"Sinmetal支2",
		"http://sinmetal2.org",
		time.Now(),
		time.Now(),
	}
	_, err = g.Put(o2)
	if err != nil {
		t.Fatal("test organization put error : %s", err.Error())
	}

	m := web.New()
	route(m)
	ts := httptest.NewServer(m)
	defer ts.Close()

	r, err := inst.NewRequest("GET", ts.URL+"/api/1/organization", nil)
	if err != nil {
		t.Error(err.Error())
	}
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Error("unexpected status code : %d, %s", w.Code, w.Body)
	}

	var o []Organization
	json.NewDecoder(w.Body).Decode(&o)
	if len(o) != 2 {
		t.Error("unexpected organization len : ", len(o))
	}
	if o[0].Id != o1.Id {
		t.Errorf("unexpected organization.id : %s, %s", o[0].Id, o1.Id)
	}
	if o[0].Name != o1.Name {
		t.Error("unexpected organization.name : %s, %s", o[0].Name, o1.Name)
	}
	if o[0].Url != o1.Url {
		t.Error("unexpected organization.url : %s, %s", o[0].Url, o1.Url)
	}
	if o[0].CreatedAt.IsZero() {
		t.Error("unexpected organization.createdAt : %s", o[0].CreatedAt)
	}

	if o[1].Id != o2.Id {
		t.Error("unexpected organization.id : %s, %s", o[1].Id, o2.Id)
	}
	if o[1].Name != o2.Name {
		t.Error("unexpected organization.name : %s, %s", o[1].Name, o2.Id)
	}
	if o[1].Url != o2.Url {
		t.Error("unexpected organization.url : ", o[1].Url, o2.Url)
	}
	if o[1].CreatedAt.IsZero() {
		t.Error("unexpected organization.createdAt : %s", o[1].CreatedAt)
	}
}

func TestPostOrganization(t *testing.T) {
	inst, c, err := aetestutil.SpinUp()
	if err != nil {
		t.Fatal(err)
	}
	defer aetestutil.SpinDown()

	g := goon.FromContext(c)

	o := &Organization{
		Id:   "sampleId",
		Name: "Sinmetal支部",
		Url:  "http://sinmetal.org",
	}
	b, err := json.Marshal(o)
	if err != nil {
		t.Fatal(err.Error())
	}

	m := web.New()
	route(m)
	ts := httptest.NewServer(m)
	defer ts.Close()

	r, err := inst.NewRequest("POST", ts.URL+"/api/1/organization", bytes.NewReader(b))
	if err != nil {
		t.Error(err.Error())
	}
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, r)
	if w.Code != http.StatusCreated {
		t.Error("unexpected status code : %d, %s", w.Code, w.Body)
	}

	var ro Organization
	json.NewDecoder(w.Body).Decode(&ro)
	if ro.Id != o.Id {
		t.Error("unexpected organization.id : ", ro.Id)
	}
	if ro.Name != o.Name {
		t.Error("unexpected organization.name : ", ro.Name)
	}
	if ro.Url != o.Url {
		t.Error("unexpected organization.url : ", o.Url)
	}
	zeroTime := time.Time{}
	if ro.CreatedAt == zeroTime {
		t.Error("unexpected organization.createdAt : ", ro.CreatedAt)
	}
	if ro.UpdatedAt == zeroTime {
		t.Error("unexpected organization.UpdatedAt : ", ro.UpdatedAt)
	}

	stored := &Organization{
		Id: o.Id,
	}
	err = g.Get(stored)
	if err != nil {
		t.Error("unexpected datastore organization, %s", err.Error())
	}
}
