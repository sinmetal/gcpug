package gcpug

import (
	"encoding/json"
	"github.com/zenazn/goji/web"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
	m := web.New()
	route(m)
	ts := httptest.NewServer(m)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/api/1/organization/sampleid")
	if err != nil {
		t.Error("unexpected")
	}
	if res.StatusCode != http.StatusOK {
		t.Error("invalid status code")
	}

	defer res.Body.Close()
	var o Organization
	json.NewDecoder(res.Body).Decode(&o)
	if o.Id != "sampleid" {
		t.Error("invalid organization.id : ", o.Id)
	}
	if o.Name != "Sinmetal支部" {
		t.Error("invalid organization.name : ", o.Name)
	}
	if o.Url != "http://sinmetal.org" {
		t.Error("invalid organization.url : ", o.Url)
	}
	zeroTime := time.Time{}
	if o.CreatedAt == zeroTime {
		t.Error("invalid organization.createdAt : ", o.CreatedAt)
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
