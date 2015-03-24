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

	"appengine"
	"github.com/mjibson/goon"

	"github.com/sinmetal/gaego_unittest_util/aetestutil"
)

type OrganizationTester struct {
}

func (t *OrganizationTester) MakeDefaultOrganization(c appengine.Context) (Organization, error) {
	g := goon.FromContext(c)

	o := Organization{
		"sampleId",
		"Sinmetal支部",
		"http://sinmetal.org",
		"http://sinmetal.org/logo.png",
		time.Now(),
		time.Now(),
	}
	_, err := g.Put(&o)
	return o, err
}

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
		"http://sinmetal.org/logo.png",
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
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code : %d, %s", w.Code, w.Body)
	}

	var ro Organization
	json.NewDecoder(w.Body).Decode(&ro)
	if ro.Id != o.Id {
		t.Fatalf("unexpected organization.id, %s != %s", ro.Id, o.Id)
	}
	if ro.Name != o.Name {
		t.Fatalf("unexpected organization.name, %s != %s", ro.Name, o.Name)
	}
	if ro.Url != o.Url {
		t.Fatalf("unexpected organization.url, %s != %s", ro.Url, o.Url)
	}
	if ro.LogoUrl != o.LogoUrl {
		t.Fatalf("unexpected organization.logoUrl, %s != %s", ro.LogoUrl, o.LogoUrl)
	}
	if ro.CreatedAt.IsZero() {
		t.Fatalf("unexpected organization.createdAt, IsZero")
	}
	if ro.UpdatedAt.IsZero() {
		t.Fatalf("unexpected organization.updatedAt, IsZero")
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
		"http://sinmetal.org/logo.png",
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
		"http://sinmetal.org/logo.png",
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
		t.Error(err)
	}
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code : %d, %s", w.Code, w.Body)
	}

	var o []Organization
	json.NewDecoder(w.Body).Decode(&o)
	if len(o) != 2 {
		t.Fatalf("unexpected organization len, %d", len(o))
	}
	if o[0].Id != o1.Id {
		t.Fatalf("unexpected organization.id, %s != %s", o[0].Id, o1.Id)
	}
	if o[0].Name != o1.Name {
		t.Fatalf("unexpected organization.name, %s != %s", o[0].Name, o1.Name)
	}
	if o[0].Url != o1.Url {
		t.Fatalf("unexpected organization.url, %s != %s", o[0].Url, o1.Url)
	}
	if o[0].LogoUrl != o1.LogoUrl {
		t.Fatalf("unexpected organization.logoUrl, %s != %s", o[0].LogoUrl, o1.LogoUrl)
	}
	if o[0].CreatedAt.IsZero() {
		t.Fatalf("unexpected organization.createdAt IsZero")
	}

	if o[1].Id != o2.Id {
		t.Fatalf("unexpected organization.id, %s != %s", o[1].Id, o2.Id)
	}
	if o[1].Name != o2.Name {
		t.Fatalf("unexpected organization.name, %s != %s", o[1].Name, o2.Id)
	}
	if o[1].Url != o2.Url {
		t.Fatalf("unexpected organization.url, %s != %s", o[1].Url, o2.Url)
	}
	if o[1].LogoUrl != o2.LogoUrl {
		t.Fatalf("unexpected organization.logoUrl, %s != %s", o[1].LogoUrl, o2.LogoUrl)
	}
	if o[1].CreatedAt.IsZero() {
		t.Fatalf("unexpected organization.createdAt, IsZero")
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
		Id:      "sampleId",
		Name:    "Sinmetal支部",
		Url:     "http://sinmetal.org",
		LogoUrl: "http://sinmetal.org/logo.png",
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
		t.Fatalf("unexpected status code : %d, %s", w.Code, w.Body)
	}

	var ro Organization
	json.NewDecoder(w.Body).Decode(&ro)
	if ro.Id != o.Id {
		t.Fatalf("unexpected organization.id, %s != %s", ro.Id, o.Id)
	}
	if ro.Name != o.Name {
		t.Fatalf("unexpected organization.name, %s != %s", ro.Name, o.Name)
	}
	if ro.Url != o.Url {
		t.Fatalf("unexpected organization.url, %s != %s", ro.Url, o.Url)
	}
	if ro.LogoUrl != o.LogoUrl {
		t.Fatalf("unexpected organization.logoUrl, %s != %s", ro.LogoUrl, o.LogoUrl)
	}
	if ro.CreatedAt.IsZero() {
		t.Fatalf("unexpected organization.createdAt, IsZero")
	}
	if ro.UpdatedAt.IsZero() {
		t.Fatalf("unexpected organization.updatedAt, IsZero")
	}

	stored := &Organization{
		Id: o.Id,
	}
	err = g.Get(stored)
	if err != nil {
		t.Fatalf("unexpected datastore organization, %s", err.Error())
	}
}

func TestPutOrganization(t *testing.T) {
	inst, c, err := aetestutil.SpinUp()
	if err != nil {
		t.Fatal(err)
	}
	defer aetestutil.SpinDown()

	g := goon.FromContext(c)

	old := &Organization{
		Id:   "sampleId",
		Name: "old支部",
	}
	_, err = g.Put(old)
	if err != nil {
		t.Fatal("test organization put error : %s", err.Error())
	}

	o := &Organization{
		Id:      "sampleId",
		Name:    "Sinmetal支部",
		Url:     "http://sinmetal.org",
		LogoUrl: "http://sinmetal.org/logo.png",
	}
	b, err := json.Marshal(o)
	if err != nil {
		t.Fatal(err.Error())
	}

	m := web.New()
	route(m)
	ts := httptest.NewServer(m)
	defer ts.Close()

	r, err := inst.NewRequest("PUT", ts.URL+"/api/1/organization", bytes.NewReader(b))
	if err != nil {
		t.Error(err.Error())
	}
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code : %d, %s", w.Code, w.Body)
	}

	var ro Organization
	json.NewDecoder(w.Body).Decode(&ro)
	if ro.Id != o.Id {
		t.Fatalf("unexpected organization.id, %s != %s", ro.Id, o.Id)
	}
	if ro.Name != o.Name {
		t.Fatalf("unexpected organization.name, %s != %s", ro.Name, o.Name)
	}
	if ro.Url != o.Url {
		t.Fatalf("unexpected organization.url, %s != %s", ro.Url, o.Url)
	}
	if ro.LogoUrl != o.LogoUrl {
		t.Fatalf("unexpected organization.logoUrl, %s != %s", ro.LogoUrl, o.LogoUrl)
	}
	if ro.CreatedAt.IsZero() {
		t.Fatalf("unexpected organization.createdAt, IsZero")
	}
	if ro.UpdatedAt.IsZero() {
		t.Fatalf("unexpected organization.updatedAt, IsZero")
	}

	stored := &Organization{
		Id: o.Id,
	}
	err = g.Get(stored)
	if err != nil {
		t.Fatalf("unexpected datastore organization, %s", err.Error())
	}

	// TODO Datastoreが更新後のデータを返してくれない
	//    b, err = json.Marshal(stored)
	//    if err != nil {
	//        t.Fatal(err.Error())
	//    }
	//    t.Fatal(string(b))
	//
	//    if stored.Id != o.Id {
	//        t.Fatalf("unexpected organization.id, %s != %s", stored.Id, o.Id)
	//    }
	//    if stored.Name != o.Name {
	//        t.Fatalf("unexpected organization.name, %s != %s", stored.Name, o.Name)
	//    }
	//    if stored.Url != o.Url {
	//        t.Fatalf("unexpected organization.url, %s != %s", stored.Url, o.Url)
	//    }
	//    if stored.LogoUrl != o.LogoUrl {
	//        t.Fatalf("unexpected organization.logoUrl, %s != %s", stored.LogoUrl, o.LogoUrl)
	//    }
	//    if stored.CreatedAt.IsZero() {
	//        t.Fatalf("unexpected organization.createdAt, IsZero")
	//    }
	//    if stored.UpdatedAt.IsZero() {
	//        t.Fatalf("unexpected organization.updatedAt, IsZero")
	//    }
}
