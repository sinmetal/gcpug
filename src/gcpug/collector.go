package gcpug

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/mjibson/goon"
	"github.com/zenazn/goji/web"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

type StackoverflowResponse struct {
	Items          []StackoverflowItem `json:"items"`
	HasMore        bool                `json:"has_more"`
	QuotaMax       int                 `json":quota_max"`
	QuotaRemaining int                 `json":quota_remaining"`
}

type StackoverflowItem struct {
	QuestionId       int       `datastore:"-" goon:"id" json:"question_id"`
	Title            string    `json:"title"`
	Link             string    `json:"link"`
	IsAnswered       bool      `json:"is_answerd"`
	ViewCount        int       `json:"view_count"`
	Score            int       `json:"score"`
	Tags             []string  `json:"tags"`
	Owner            Owner     `json:"owner"`
	CreationDate     EpochTime `json:"creation_date"`
	LastActivityDate EpochTime `json:"last_activity_date"`
}

type Stackoverflow struct {
	QuestionId       int64     `datastore:"-" goon:"id" json:"questionID"`
	Title            string    `json:"title" datastore:",noindex`
	Link             string    `json:"link" datastore:",noindex`
	IsAnswered       bool      `json:"isAnswerd" datastore:",noindex`
	ViewCount        int       `json:"viewCount" datastore:",noindex`
	Score            int       `json:"score" datastore:",noindex`
	Tags             []string  `json:"tags" datastore:",noindex`
	Owner            Owner     `json:"owner"`
	CreationDate     time.Time `json:"creationDate" datastore:",noindex`
	LastActivityDate time.Time `json:"lastActivityDate" datastore:",noindex`
}

type Owner struct {
	Requtation   int    `json:"reputation" datastore:",noindex`
	UserId       int    `json:"user_id" datastore:",noindex`
	UserType     string `json:"user_type" datastore:",noindex`
	ProfileImage string `json:"profile_image" datastore:",noindex`
	DisplayName  string `json:"display_name" datastore:",noindex`
	Link         string `json:"link" datastore:",noindex`
}

type EpochTime time.Time

func (t *EpochTime) UnmarshalJSON(buf []byte) error {
	value, err := strconv.ParseInt(string(buf), 10, 64)
	if err != nil {
		return err
	}
	*t = EpochTime(time.Unix(value, 0))
	return nil
}

func (t *EpochTime) MarshalJSON() ([]byte, error) {
	epoch := (time.Time)(*t).Unix()
	bs := []byte(strconv.FormatInt(epoch, 10))
	return bs, nil
}

type CollectorApi struct {
	Config PugConfig
}

func SetUpCollector(m *web.Mux) {
	api := CollectorApi{}

	m.Get("/cron/1/collector", api.Get)
}

func (a *CollectorApi) Get(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	appId := appengine.AppID(ctx)
	if appengine.IsDevAppServer() == false && appId != "gcp-ug" {
		log.Infof(ctx, "do nothing. only run gcp-ug. appId = %s", appId)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		return
	}

	g := goon.NewGoon(r)
	s := &PugConfigService{}
	pc, err := s.Get(g)
	if err != nil {
		er := ErrorResponse{
			http.StatusInternalServerError,
			[]string{err.Error()},
		}
		log.Errorf(ctx, fmt.Sprintf("pug config get error : ", err.Error()))
		er.Write(w)
		return
	}
	a.Config = pc

	err = a.PullStackoverflow(ctx, "google-app-engine")
	if err != nil {
		er := ErrorResponse{
			http.StatusInternalServerError,
			[]string{err.Error()},
		}
		log.Errorf(ctx, fmt.Sprintf("pull stackoverflow error : ", err.Error()))
		er.Write(w)
		return
	}

	a.PullStackoverflow(ctx, "google-compute-engine")
	if err != nil {
		er := ErrorResponse{
			http.StatusInternalServerError,
			[]string{err.Error()},
		}
		log.Errorf(ctx, fmt.Sprintf("pull stackoverflow error : ", err.Error()))
		er.Write(w)
		return
	}

	a.PullStackoverflow(ctx, "google-cloud-sql")
	if err != nil {
		er := ErrorResponse{
			http.StatusInternalServerError,
			[]string{err.Error()},
		}
		log.Errorf(ctx, fmt.Sprintf("pull stackoverflow error : ", err.Error()))
		er.Write(w)
		return
	}

	a.PullStackoverflow(ctx, "google-bigquery")
	if err != nil {
		er := ErrorResponse{
			http.StatusInternalServerError,
			[]string{err.Error()},
		}
		log.Errorf(ctx, fmt.Sprintf("pull stackoverflow error : ", err.Error()))
		er.Write(w)
		return
	}

	a.PullStackoverflow(ctx, "google-cloud-storage")
	if err != nil {
		er := ErrorResponse{
			http.StatusInternalServerError,
			[]string{err.Error()},
		}
		log.Errorf(ctx, fmt.Sprintf("pull stackoverflow error : ", err.Error()))
		er.Write(w)
		return
	}

	a.PullStackoverflow(ctx, "google-cloud-datastore")
	if err != nil {
		er := ErrorResponse{
			http.StatusInternalServerError,
			[]string{err.Error()},
		}
		log.Errorf(ctx, fmt.Sprintf("pull stackoverflow error : ", err.Error()))
		er.Write(w)
		return
	}

	a.PullStackoverflow(ctx, "google-cloud-endpoints")
	if err != nil {
		er := ErrorResponse{
			http.StatusInternalServerError,
			[]string{err.Error()},
		}
		log.Errorf(ctx, fmt.Sprintf("pull stackoverflow error : ", err.Error()))
		er.Write(w)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (a *CollectorApi) PullStackoverflow(ctx context.Context, tag string) error {
	log.Infof(ctx, "fetch %s start.", tag)

	client := urlfetch.Client(ctx)
	uri := fmt.Sprintf("https://api.stackexchange.com/2.2/questions?order=desc&sort=activity&tagged=%s&site=ja.stackoverflow", tag)
	resp, err := client.Get(uri)
	if err != nil {
		log.Errorf(ctx, "%s", err.Error())
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	log.Infof(ctx, "%s", string(body))
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("stackoverflow error : code = %d", resp.StatusCode))
	}

	stackArray, err := a.ParseJson(body)
	if err != nil {
		log.Errorf(ctx, "%s", err.Error())
		return err
	}
	log.Infof(ctx, "size = %d", len(stackArray))

	for _, stack := range stackArray {
		key := datastore.NewKey(ctx, "Stackoverflow", "", int64(stack.QuestionId), nil)
		err = datastore.RunInTransaction(ctx, func(ctx context.Context) error {
			stored := &Stackoverflow{}
			err := datastore.Get(ctx, key, stored)
			if err == nil {
				return nil
			}
			if err != datastore.ErrNoSuchEntity {
				log.Errorf(ctx, "%s", err.Error())
				return err
			}

			stored.Title = stack.Title
			stored.Link = stack.Link
			stored.IsAnswered = stack.IsAnswered
			stored.ViewCount = stack.ViewCount
			stored.Score = stack.Score
			stored.Tags = stack.Tags
			stored.Owner = stack.Owner
			stored.CreationDate = (time.Time)(stack.CreationDate)
			stored.LastActivityDate = (time.Time)(stack.LastActivityDate)
			_, err = datastore.Put(ctx, key, stored)
			if err != nil {
				return err
			}

			sm := SlackMessage{}
			sm.Set(stored)
			_, err = a.PostToSlack(ctx, sm)
			if err != nil {
				return err
			}

			return nil
		}, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *CollectorApi) ParseJson(body []byte) ([]StackoverflowItem, error) {
	var stackRes StackoverflowResponse
	err := json.Unmarshal(body, &stackRes)
	if err != nil {
		fmt.Printf("json parse error: %v", err)
		return nil, err
	}
	return stackRes.Items, nil
}

type SlackMessage struct {
	UserName    string            `json:"username"`
	IconUrl     string            `json:"icon_url"`
	Text        string            `json:"text"`
	Attachments []SlackAttachment `json:"attachments"`
}

type SlackAttachment struct {
	Color      string       `json:"color"`
	AuthorName string       `json:"author_name"`
	AuthorLink string       `json:"author_link"`
	AuthorIcon string       `json:"author_icon"`
	Title      string       `json:"title"`
	TitleLink  string       `json:"title_link"`
	Fields     []SlackField `json:"fields"`
}

type SlackField struct {
	Title string `json:"title"`
}

func (sm *SlackMessage) Set(s *Stackoverflow) {
	fields := make([]SlackField, 0)
	for _, title := range s.Tags {
		fields = append(fields, SlackField{
			Title: title,
		})
	}

	sa := SlackAttachment{
		Color:      "#36a64f",
		AuthorName: s.Owner.DisplayName,
		AuthorLink: s.Owner.Link,
		AuthorIcon: s.Owner.ProfileImage,
		Title:      s.Title,
		TitleLink:  s.Link,
		Fields:     fields,
	}

	sm.UserName = "gcpug"
	sm.IconUrl = "http://gcpug.jp/images/logo_box.png"
	sm.Text = s.Title
	sm.Attachments = []SlackAttachment{sa}
}

func (a *CollectorApi) PostToSlack(c context.Context, message SlackMessage) (resp *http.Response, err error) {
	client := urlfetch.Client(c)

	body, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))
	return client.Post(
		a.Config.SlackPostUrl,
		"application/json",
		bytes.NewReader(body))
}
