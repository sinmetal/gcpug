package gcpug

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/zenazn/goji/web"

	"github.com/mjibson/goon"

	"github.com/sinmetal/gaego_unittest_util/aetestutil"
)

func TestPostPugEvent(t *testing.T) {
	inst, c, err := aetestutil.SpinUp()
	if err != nil {
		t.Fatal(err)
	}
	defer aetestutil.SpinDown()

	g := goon.FromContext(c)

	o := &Organization{
		Id: "organizationId",
	}

	oKey := g.Key(o)

	pe := &PugEvent{
		OrganizationKey: oKey,
		Title:           "GAEハンズオン",
		Url:             "http://example.com",
		StartAt:         time.Now(),
	}

	b, err := json.Marshal(pe)
	if err != nil {
		t.Fatal(err)
	}

	m := web.New()
	route(m)
	ts := httptest.NewServer(m)
	defer ts.Close()

	r, err := inst.NewRequest("POST", ts.URL+"/api/1/event", bytes.NewReader(b))
	if err != nil {
		t.Error(err.Error())
	}
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, r)
	if w.Code != http.StatusCreated {
		t.Fatalf("unexpected status code : %d, %s", w.Code, w.Body)
	}

	var re PugEvent
	err = json.NewDecoder(w.Body).Decode(&re)
	if err != nil {
		t.Fatal(err)
	}
	if re.Id == "" {
		t.Fatalf("unexpected pug event id, empty")
	}
	if re.Title != pe.Title {
		t.Fatalf("unexpected pug event title, %s != %s", re.Title, pe.Title)
	}
	if *re.OrganizationKey != *pe.OrganizationKey {
		t.Fatalf("unexpected pug event organization key, %v != %v", re.OrganizationKey, pe.OrganizationKey)
	}
	if re.StartAt != pe.StartAt {
		t.Fatalf("unexpected pug envet start at, %s != %s", re.StartAt, pe.StartAt)
	}
	if re.CreatedAt.IsZero() {
		t.Fatalf("unexpected pug event created at, IsZero")
	}
	if re.UpdatedAt.IsZero() {
		t.Fatalf("unexpected pug event updated at, IsZero")
	}

	stored := &PugEvent{
		Id: re.Id,
	}
	err = g.Get(stored)
	if err != nil {
		t.Fatalf("unexpected datastore pug event, %s", err.Error())
	}
}

func TestPugEventSave(t *testing.T) {
	_, c, err := aetestutil.SpinUp()
	if err != nil {
		t.Fatal(err)
	}
	defer aetestutil.SpinDown()

	g := goon.FromContext(c)

	o := &Organization{
		Id: "organizationId",
	}
	key := g.Key(o)

	pe := &PugEvent{
		Id:              "hogeId",
		OrganizationKey: key,
	}

	_, err = g.Put(pe)
	if err != nil {
		t.Fatal(err)
	}
	if pe.CreatedAt.IsZero() {
		t.Fatalf("unexpected createdAt. createdAt is zero value")
	}
	if pe.UpdatedAt.IsZero() {
		t.Fatalf("unexpected updatedAt. updatedAt is zero value")
	}

	var after PugEvent
	peJson := `{"Id":"hogeId","OrganizationKey":"agxkZXZ-dW5pdHRlc3RyIAsSDE9yZ2FuaXphdGlvbiIOb3JnYW5pemF0aW9uSWQM","Title":"hogeTitle","Url":"http://example.com","StartAt":"2015-03-09T19:47:16.801665955+09:00","CreatedAt":"2015-03-09T19:47:16.801665955+09:00","UpdatedAt":"2015-03-09T19:47:16.801665955+09:00"}`
	err = json.Unmarshal([]byte(peJson), &after)
	if err != nil {
		t.Error(err)
	}
	if after.Id != "hogeId" {
		t.Fatalf("unexpected id. id = %s")
	}
	if *after.OrganizationKey != *key {
		t.Fatalf("unexpected organization key : %s != %s", after.OrganizationKey, key)
	}

	expectedStartAt, err := time.Parse(
		"2006-01-02T15:04:05.999999999-07:00", // スキャンフォーマット
		"2015-03-09T19:47:16.801665955+09:00") // パースしたい文字列
	if err != nil {
		t.Fatal(err)
	}
	if after.StartAt != expectedStartAt {
		t.Fatalf("unexpected startAt. %s != %s", after.StartAt.String(), expectedStartAt)
	}
}
