package gcpug

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/zenazn/goji/web"

	"github.com/mjibson/goon"
	"google.golang.org/appengine/aetest"
)

type PugEventTester struct {
}

func (t *PugEventTester) MakePugEvent(r *http.Request, o Organization, n int) (PugEvent, error) {
	g := goon.NewGoon(r)

	pe := PugEvent{
		Id:             fmt.Sprintf("test%d", n),
		OrganizationId: o.Id,
		Title:          fmt.Sprintf("GAEハンズオン%d", n),
		Url:            fmt.Sprintf("http://example%d.com", n),
		StartAt:        time.Now(),
	}
	_, err := g.Put(&pe)
	return pe, err
}

func (pet *PugEventTester) Equal(t *testing.T, pe1 PugEvent, pe2 PugEvent) {
	if pe1.Id != pe2.Id {
		t.Fatalf("unexpected response pugEvent.Id, %s != %s", pe1.Id, pe2.Id)
	}
	if pe1.OrganizationId != pe2.OrganizationId {
		t.Fatalf("unexpected response pugEvent.OrganizationId, %v != %v", pe1.OrganizationId, pe2.OrganizationId)
	}
	if pe1.Title != pe2.Title {
		t.Fatalf("unexpected response pugEvent.Title, %s != %s", pe1.Title, pe2.Title)
	}
	if pe1.Description != pe2.Description {
		t.Fatalf("unexpected response pugEvent description, %s != %s", pe1.Description, pe2.Description)
	}
	if pe1.Url != pe2.Url {
		t.Fatalf("unexpected response pugEvent.Url, %s != %s", pe1.Url, pe2.Url)
	}
	if EqualYYYYMMDDHHMMSS(pe1.StartAt, pe2.StartAt) == false {
		t.Fatalf("unexpected response pugEvent.StartAt, %s != %s", pe1.StartAt, pe2.StartAt)
	}
	if EqualYYYYMMDDHHMMSS(pe1.CreatedAt, pe2.CreatedAt) == false {
		t.Fatalf("unexpected response pugEvent.CreatedAt, %s != %s", pe1.CreatedAt, pe2.CreatedAt)
	}
	if EqualYYYYMMDDHHMMSS(pe1.UpdatedAt, pe2.UpdatedAt) == false {
		t.Fatalf("unexpected response pugEvent.UpdatedAt, %s != %s", pe1.UpdatedAt, pe2.UpdatedAt)
	}
}

func EqualYYYYMMDDHHMMSS(t1 time.Time, t2 time.Time) bool {
	const format = "2006-01-02T15:04:05-07:00"
	return t1.Format(format) == t2.Format(format)
}

func TestPostPugEvent(t *testing.T) {
	t.Parallel()

	opt := &aetest.Options{AppID: "unittest", StronglyConsistentDatastore: true}
	inst, err := aetest.NewInstance(opt)
	defer inst.Close()
	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("fatal new request error : %s", err.Error())
	}

	g := goon.NewGoon(req)

	o := &Organization{
		Id: "organizationId",
	}

	pe := &PugEvent{
		OrganizationId: o.Id,
		Title:          "GAEハンズオン",
		Description:    "初心者のためのGAEハンズオン！",
		Url:            "http://example.com",
		StartAt:        time.Now(),
	}

	b, err := json.Marshal(pe)
	if err != nil {
		t.Fatal(err)
	}

	m := web.New()
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
	if re.Description != pe.Description {
		t.Fatalf("unexpected pug event description, %s != %s", re.Description, pe.Description)
	}
	if re.OrganizationId != pe.OrganizationId {
		t.Fatalf("unexpected pug event organization id, %v != %v", re.OrganizationId, pe.OrganizationId)
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

func TestPutPugEventBadRequest(t *testing.T) {
	t.Parallel()

	opt := &aetest.Options{AppID: "unittest", StronglyConsistentDatastore: true}
	inst, err := aetest.NewInstance(opt)
	defer inst.Close()
	_, err = inst.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("fatal new request error : %s", err.Error())
	}

	o := &Organization{
		Id: "organizationId",
	}

	pe := &PugEvent{
		OrganizationId: o.Id,
		Title:          "GAEハンズオン",
		Description:    "初心者のためのGAEハンズオン！",
		Url:            "http://example.com",
		StartAt:        time.Now(),
	}

	b, err := json.Marshal(pe)
	if err != nil {
		t.Fatal(err)
	}

	m := web.New()
	ts := httptest.NewServer(m)
	defer ts.Close()

	r, err := inst.NewRequest("PUT", ts.URL+"/api/1/event", bytes.NewReader(b))
	if err != nil {
		t.Error(err.Error())
	}
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status code : %d, %s", w.Code, w.Body)
	}

	var messages []string
	err = json.NewDecoder(w.Body).Decode(&messages)
	if err != nil {
		t.Fatalf("unexpected response body")
	}
	if messages[0] != "id is required." {
		t.Fatalf("unexpected response body")
	}
}

func TestPutPugEventNothing(t *testing.T) {
	t.Parallel()

	opt := &aetest.Options{AppID: "unittest", StronglyConsistentDatastore: true}
	inst, err := aetest.NewInstance(opt)
	defer inst.Close()
	_, err = inst.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("fatal new request error : %s", err.Error())
	}

	o := &Organization{
		Id: "organizationId",
	}

	pe := &PugEvent{
		Id:             "nothingid",
		OrganizationId: o.Id,
		Title:          "GAEハンズオン",
		Description:    "初心者のためのGAEハンズオン！",
		Url:            "http://example.com",
		StartAt:        time.Now(),
	}

	b, err := json.Marshal(pe)
	if err != nil {
		t.Fatal(err)
	}

	m := web.New()
	ts := httptest.NewServer(m)
	defer ts.Close()

	r, err := inst.NewRequest("PUT", ts.URL+"/api/1/event", bytes.NewReader(b))
	if err != nil {
		t.Error(err.Error())
	}
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("unexpected status code : %d, %s", w.Code, w.Body)
	}

	var messages []string
	err = json.NewDecoder(w.Body).Decode(&messages)
	if err != nil {
		t.Fatalf("unexpected response body : %s")
	}
	if messages[0] != "datastore: no such entity" {
		t.Fatalf("unexpected response body")
	}
}

func TestPutPugEvent(t *testing.T) {
	t.Parallel()

	opt := &aetest.Options{AppID: "unittest", StronglyConsistentDatastore: true}
	inst, err := aetest.NewInstance(opt)
	defer inst.Close()
	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("fatal new request error : %s", err.Error())
	}

	g := goon.NewGoon(req)

	o := &Organization{
		Id: "organizationId",
	}

	data := &PugEvent{
		Id:             "testevent",
		OrganizationId: o.Id,
		Title:          "GAEハンズオン",
		Description:    "初心者のためのGAEハンズオン！",
		Url:            "http://example.com",
		StartAt:        time.Now(),
	}
	err = data.Create(g)
	if err != nil {
		t.Fatal(err)
	}

	pe := &PugEvent{
		Id:             data.Id,
		OrganizationId: o.Id,
		Title:          "GAEハンズオン-new",
		Description:    "上級者のためのGAE",
		Url:            "https://example.com?gae=gae",
		StartAt:        time.Now(),
	}

	b, err := json.Marshal(pe)
	if err != nil {
		t.Fatal(err)
	}

	m := web.New()
	ts := httptest.NewServer(m)
	defer ts.Close()

	r, err := inst.NewRequest("PUT", ts.URL+"/api/1/event", bytes.NewReader(b))
	if err != nil {
		t.Error(err.Error())
	}
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code : %d, %s", w.Code, w.Body)
	}

	var re PugEvent
	err = json.NewDecoder(w.Body).Decode(&re)
	if err != nil {
		t.Fatal(err)
	}
	if re.Id != pe.Id {
		t.Fatalf("unexpected pug event id, %s != %s", re.Id, pe.Id)
	}
	if re.Title != pe.Title {
		t.Fatalf("unexpected pug event title, %s != %s", re.Title, pe.Title)
	}
	if re.Description != pe.Description {
		t.Fatalf("unexpected pug event description, %s != %s", re.Description, pe.Description)
	}
	if re.OrganizationId != pe.OrganizationId {
		t.Fatalf("unexpected pug event organization id, %v != %v", re.OrganizationId, pe.OrganizationId)
	}
	if re.StartAt != pe.StartAt {
		t.Fatalf("unexpected pug envet start at, %s != %s", re.StartAt, pe.StartAt)
	}
	if re.CreatedAt == data.CreatedAt {
		t.Fatalf("unexpected pug event created at, %s != %s", re.CreatedAt, pe.CreatedAt)
	}
	if re.UpdatedAt.Unix() > data.UpdatedAt.Unix() {
		t.Fatalf("unexpected pug event updated at, %s != %s", re.UpdatedAt, pe.UpdatedAt)
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
	t.Parallel()

	opt := &aetest.Options{AppID: "unittest", StronglyConsistentDatastore: true}
	inst, err := aetest.NewInstance(opt)
	defer inst.Close()
	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("fatal new request error : %s", err.Error())
	}

	g := goon.NewGoon(req)

	o := &Organization{
		Id: "organizationId",
	}

	pe := &PugEvent{
		Id:             "hogeId",
		OrganizationId: o.Id,
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
	peJson := `{"Id":"hogeId","OrganizationId":"organizationId","Title":"hogeTitle","Url":"http://example.com","StartAt":"2015-03-09T19:47:16.801665955+09:00","CreatedAt":"2015-03-09T19:47:16.801665955+09:00","UpdatedAt":"2015-03-09T19:47:16.801665955+09:00"}`
	err = json.Unmarshal([]byte(peJson), &after)
	if err != nil {
		t.Error(err)
	}
	if after.Id != "hogeId" {
		t.Fatalf("unexpected id. id = %s")
	}
	if after.OrganizationId != o.Id {
		t.Fatalf("unexpected organization id : %s != %s", after.OrganizationId, o.Id)
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

func TestListPugEvent(t *testing.T) {
	t.Parallel()

	opt := &aetest.Options{AppID: "unittest", StronglyConsistentDatastore: true}
	inst, err := aetest.NewInstance(opt)
	defer inst.Close()
	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("fatal new request error : %s", err.Error())
	}

	ot := OrganizationTester{}
	o, err := ot.MakeDefaultOrganization(req)
	if err != nil {
		t.Error(err)
	}

	pet := PugEventTester{}
	pe1, err := pet.MakePugEvent(req, o, 1)
	if err != nil {
		t.Error(err)
	}

	pe2, err := pet.MakePugEvent(req, o, 2)
	if err != nil {
		t.Error(err)
	}

	m := web.New()
	ts := httptest.NewServer(m)
	defer ts.Close()

	r, err := inst.NewRequest("GET", ts.URL+"/api/1/event", nil)
	if err != nil {
		t.Error(err.Error())
	}
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code : %d, %s", w.Code, w.Body)
	}

	var pes []PugEvent
	json.NewDecoder(w.Body).Decode(&pes)
	if len(pes) != 2 {
		t.Fatalf("unexpected response pugEvent length, %d", len(pes))
	}
	pet.Equal(t, pes[0], pe2)
	pet.Equal(t, pes[1], pe1)
}
