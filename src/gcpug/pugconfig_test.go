package gcpug

import (
	"bytes"
	"encoding/json"
	"github.com/zenazn/goji/web"
	"net/http"
	"net/http/httptest"
	"testing"

	"appengine"
	"appengine/aetest"

	"github.com/mjibson/goon"
)

type PugConfigTester struct {
}

func TestPugConfigPost(t *testing.T) {
	t.Parallel()

	opt := &aetest.Options{AppID: "unittest", StronglyConsistentDatastore: true}
	inst, err := aetest.NewInstance(opt)
	defer inst.Close()
	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("fatal new request error : %s", err.Error())
	}

	c := appengine.NewContext(req)

	g := goon.FromContext(c)

	con := &PugConfig{
		ClientId:     "hoge-clinet-id",
		ClientSecret: "hoge-client-secret",
		SlackPostUrl: "http://example.com",
	}
	b, err := json.Marshal(con)
	if err != nil {
		t.Fatal(err)
	}

	m := web.New()
	ts := httptest.NewServer(m)
	defer ts.Close()

	r, err := inst.NewRequest("POST", ts.URL+"/admin/api/1/config", bytes.NewReader(b))
	if err != nil {
		t.Error(err.Error())
	}
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code : %d, %s", w.Code, w.Body)
	}

	var rpc PugConfig
	json.NewDecoder(w.Body).Decode(&rpc)
	if rpc.Id != pugConfigId {
		t.Fatalf("unexpected pugConfig.id, %s != %s", rpc.Id, pugConfigId)
	}
	if rpc.ClientId != con.ClientId {
		t.Fatalf("unexpected pugConfig.ClinetId, %s != %s", rpc.ClientId, con.ClientId)
	}
	if rpc.ClientSecret != con.ClientSecret {
		t.Fatalf("unexpected pugConfig.ClinetSecret, %s != %s", rpc.ClientSecret, con.ClientSecret)
	}
	if rpc.SlackPostUrl != con.SlackPostUrl {
		t.Fatalf("unexpceted pugConfig.SlackPostUrl, %s != %s", rpc.SlackPostUrl, con.SlackPostUrl)
	}
	if rpc.CreatedAt.IsZero() {
		t.Fatalf("unexpected pugConfig.createdAt, IsZero")
	}
	if rpc.UpdatedAt.IsZero() {
		t.Fatalf("unexpected pugConfig.updatedAt, IsZero")
	}

	stored := &PugConfig{
		Id: pugConfigId,
	}
	err = g.Get(stored)
	if err != nil {
		t.Fatalf("unexpected datastore pugConfig, %v", err)
	}
}

func TestPugConfigGet(t *testing.T) {
	t.Parallel()

	opt := &aetest.Options{AppID: "unittest", StronglyConsistentDatastore: true}
	inst, err := aetest.NewInstance(opt)
	defer inst.Close()
	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("fatal new request error : %s", err.Error())
	}

	c := appengine.NewContext(req)

	g := goon.FromContext(c)

	con := &PugConfig{
		Id:           pugConfigId,
		ClientId:     "hoge-clinet-id",
		ClientSecret: "hoge-client-secret",
		SlackPostUrl: "http://example.com",
	}
	_, err = g.Put(con)
	if err != nil {
		t.Error(err)
	}

	s := &PugConfigService{}
	pc, err := s.Get(c)
	if err != nil {
		t.Error(err)
	}
	if pc.Id != pugConfigId {
		t.Fatalf("unexpected pugConfig.id, %s != %s", pc.Id, pugConfigId)
	}
	if pc.ClientId != con.ClientId {
		t.Fatalf("unexpected pugConfig.ClinetId, %s != %s", pc.ClientId, con.ClientId)
	}
	if pc.ClientSecret != con.ClientSecret {
		t.Fatalf("unexpected pugConfig.ClinetSecret, %s != %s", pc.ClientSecret, con.ClientSecret)
	}
	if pc.SlackPostUrl != con.SlackPostUrl {
		t.Fatalf("unexpected pugConfig.SlackPostUrl, %s != %s", pc.SlackPostUrl, con.SlackPostUrl)
	}
	if pc.CreatedAt.IsZero() {
		t.Fatalf("unexpected pugConfig.createdAt, IsZero")
	}
	if pc.UpdatedAt.IsZero() {
		t.Fatalf("unexpected pugConfig.updatedAt, IsZero")
	}
}
