package gcpug

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"github.com/mjibson/goon"
	gengine "google.golang.org/appengine"
	"google.golang.org/appengine/log"
	gmemcache "google.golang.org/appengine/memcache"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gauth "google.golang.org/api/oauth2/v2"

	plus "google.golang.org/api/plus/v1"
)

type requestParam struct {
	Host       string
	Method     string
	UrlHost    string
	Fragment   string
	Path       string
	Scheme     string
	Opaque     string
	RawQuery   string
	RemoteAddr string
	RequestURI string
	UserAgent  string
}

const (
	randStateForAuthToMemcacheKey = "randStateForAuthToMemcacheKey"
)

type Member struct {
	Email       string    `datastore:"-" goon:"id" json:"email"`    // Email
	Id          string    `json:"-"`                                // Google Id
	NickName    string    `json:"nickName"`                         // ニックネーム
	Name        string    `json:"name"`                             // G+ 名前
	FamilyName  string    `json:"familyName" datastore:",noindex"`  // G+ 姓
	GivenName   string    `json:"givenName" datastore:",noindex"`   // G+ 名
	PlusLink    string    `json:"plusLink" datastore:",noindex"`    // G+ Profile Link
	PictureLink string    `json:"pictureLink" datastore:",noindex"` // G+ Picture Link
	GithubId    string    `json:"githubId" datastore:",noindex"`    // Github Id
	QiitaId     string    `json:"qiitaId" datastore:",noindex"`     // Qiita Id
	TwitterId   string    `json:"twitterId" datastore:",noindex"`   // Twitter Id
	FacebookId  string    `json:"facebookId" datastore:",noindex"`  // Faebook Id
	BlogLink    string    `json:"blogLink" datastore:",noindex"`    // Blog Link
	CreatedAt   time.Time `json:"createdAt"`                        // 作成日時
	UpdatedAt   time.Time `json:"updatedAt"`                        // 更新日時
}

type MemberApi struct {
}

func SetUpMember(m *web.Mux) {
	api := MemberApi{}

	m.Get("/api/1/login", api.Login)
	m.Get("/oauth2callback", api.OAuth2Callback)
}

func (a *MemberApi) getConfig(r *http.Request) (*oauth2.Config, error) {
	ac := appengine.NewContext(r)

	pcs := &PugConfigService{}
	config, err := pcs.Get(ac)
	if err != nil {
		return &oauth2.Config{}, err
	}

	protocol := "https"
	if appengine.IsDevAppServer() {
		protocol = "http"
	}
	redirectUrl := fmt.Sprintf("%s://%s/oauth2callback", protocol, r.Host)

	return &oauth2.Config{
		ClientID:     config.ClientId,
		ClientSecret: config.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectUrl,
		Scopes:       []string{plus.PlusMeScope, plus.UserinfoEmailScope},
	}, nil
}

func (a *MemberApi) Login(c web.C, w http.ResponseWriter, r *http.Request) {
	ac := appengine.NewContext(r)

	config, err := a.getConfig(r)
	if err != nil {
		ac.Errorf("pug config get error, %v", err)
		er := ErrorResponse{
			http.StatusInternalServerError,
			[]string{"config get error"},
		}
		er.Write(w)
		return
	}

	randState := fmt.Sprintf("st%d", time.Now().UnixNano())
	authUrl := config.AuthCodeURL(randState)
	ac.Infof("auth url = %s", authUrl)

	item := &memcache.Item{
		Key:        fmt.Sprintf("%s-_-%s", randStateForAuthToMemcacheKey, randState),
		Value:      []byte(randState),
		Expiration: 3 * time.Minute,
	}
	memcache.Add(ac, item)

	http.Redirect(w, r, authUrl, http.StatusFound)
}

func (a *MemberApi) OAuth2Callback(c web.C, w http.ResponseWriter, r *http.Request) {
	ac := gengine.NewContext(r)

	p := &requestParam{
		r.Host,
		r.Method,
		r.URL.Host,
		r.URL.Fragment,
		r.URL.Path,
		r.URL.Scheme,
		r.URL.Opaque,
		r.URL.RawQuery,
		r.RemoteAddr,
		r.RequestURI,
		r.UserAgent(),
	}

	_, err := json.Marshal(p)
	if err != nil {
		log.Errorf(ac, "handler error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	config, err := a.getConfig(r)
	if err != nil {
		log.Errorf(ac, "pug config get error, %v", err)
		er := ErrorResponse{
			http.StatusInternalServerError,
			[]string{"config get error"},
		}
		er.Write(w)
		return
	}

	stateMemKey := fmt.Sprintf("%s-_-%s", randStateForAuthToMemcacheKey, r.FormValue("state"))
	item, err := gmemcache.Get(ac, stateMemKey)
	if err != nil {
		log.Errorf(ac, "memcache get error, %v", err)
		er := ErrorResponse{
			http.StatusUnauthorized,
			[]string{"unauthorized"},
		}
		er.Write(w)
		return
	}

	if r.FormValue("state") != string(item.Value) {
		log.Warningf(ac, "State doesn't match: req = %#v", "")
		er := ErrorResponse{
			http.StatusUnauthorized,
			[]string{"unauthorized"},
		}
		er.Write(w)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		log.Errorf(ac, "token not found.")
	}
	token, err := config.Exchange(ac, code)
	if err != nil {
		log.Errorf(ac, "Token exchange error: %v", err)
	}
	_, err = json.Marshal(&token)
	if err != nil {
		log.Errorf(ac, "token json marshal error: %v", err)
	}

	oauthClient := config.Client(ac, token)
	s, err := gauth.New(oauthClient)
	if err != nil {
		log.Errorf(ac, "gauth new error: %v", err)
	}
	me := gauth.NewUserinfoV2MeService(s)
	ui, err := me.Get().Do()
	if err != nil {
		log.Errorf(ac, "get user info from plus error: %v", err)
	}

	m := &Member{
		Email:       ui.Email,
		Id:          ui.Id,
		Name:        ui.Name,
		FamilyName:  ui.FamilyName,
		GivenName:   ui.GivenName,
		PlusLink:    ui.Link,
		PictureLink: ui.Picture,
	}
	g := goon.NewGoon(r)
	err = m.PutByLogin(g)
	if err != nil {
		log.Errorf(ac, "member put error, %v", err)
		er := ErrorResponse{
			http.StatusInternalServerError,
			[]string{"member put error"},
		}
		er.Write(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(m)
}

func (m *Member) PutByLogin(g *goon.Goon) error {
	return g.RunInTransaction(func(g *goon.Goon) error {
		stored := &Member{
			Email: m.Email,
		}
		err := g.Get(stored)
		if err == nil {
			stored.Name = m.Name
			stored.FamilyName = m.FamilyName
			stored.GivenName = m.GivenName
			stored.PictureLink = m.PictureLink
			stored.PlusLink = m.PlusLink

			_, err = g.Put(stored)
			if err != nil {
				return err
			}
			*m = *stored

			return nil
		} else if err == datastore.ErrNoSuchEntity {
			_, err = g.Put(m)
			if err != nil {
				return err
			}

			return nil
		} else {
			return err
		}
	}, nil)
}

func (m *Member) Load(c <-chan datastore.Property) error {
	if err := datastore.LoadStruct(m, c); err != nil {
		return err
	}

	return nil
}

func (m *Member) Save(c chan<- datastore.Property) error {
	now := time.Now()
	m.UpdatedAt = now

	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}

	if err := datastore.SaveStruct(m, c); err != nil {
		return err
	}
	return nil
}
