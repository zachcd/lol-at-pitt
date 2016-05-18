package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/lab-d8/lol-at-pitt/ols"
	"github.com/lab-d8/oauth2"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/rs/cors"
	goauth2 "golang.org/x/oauth2"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"net/url"
	"time"
)

func InitMiddleware(m *martini.ClassicMartini) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	m.Handlers(PARAMS,
		DB(),
		sessions.Sessions("lol_session", sessions.NewCookieStore([]byte("secret123"))),
		oauth2.Facebook(
			&goauth2.Config{
				ClientID:     ClientId,
				ClientSecret: ApiSecret,
				Scopes:       []string{"public_profile", "email", "user_friends"},
				RedirectURL:  "http://www.lol-at-pitt.com/oauth2callback",
			},
		),
		c.HandlerFunc,
		render.Renderer(render.Options{Directory: TemplatesLocation}),
		martini.Static("resources/public", martini.StaticOptions{Prefix: "/public"}),
	)
}

func InitDebugMiddleware(m *martini.ClassicMartini) {
	m.Use(PARAMS)
	m.Use(DB())
	m.Use(sessions.Sessions("lol_session", sessions.NewCookieStore([]byte("secret123"))))
	m.Use(render.Renderer(render.Options{Directory: TemplatesLocation}))
	m.Use(martini.Static("resources/public", martini.StaticOptions{Prefix: "/public"}))
	SetId("1", "10153410152015744", "Sean Myers") // Me. Set these to act like facebook, using a nice cache
}

// PARAMS is a middleware binder for injecting the params into each handler
func PARAMS(req *http.Request, c martini.Context) {
	req.ParseForm()
	response := req.Form
	c.Map(response)
	c.Next()
}

// DB is a middleware binder that injects the mongo db into each handler
func DB() martini.Handler {
	session, err := mgo.Dial(ols.MongoLocation)
	if err != nil {
		panic(err)
	}

	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB(ols.DatabaseName))
		defer s.Close()
		c.Next()
	}
}

func Permissions(permissionName string) martini.Handler {
	return func(token oauth2.Tokens, w http.ResponseWriter, r *http.Request, c martini.Context) {
		if token == nil || token.Expired() {
			next := url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, oauth2.PathLogin+"?next="+next, 302)
			return
		}
		id, err := GetId(token.Access())
		if err != nil {
			log.Printf("Error getting player token id:", err.Error())
			http.Redirect(w, r, "/error", 302)
			return
		}
		user := ols.GetUserDAO().GetUserFB(id)
		if user.LeagueId == 0 {
			next := url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, "/register?next="+next, 302)
		}

		// TODO - fix this
		if !true {
			http.Redirect(w, r, "/error", 302)
		}
		c.Map(user)
		c.Next()

	}
}

var CaptainRequiredFunc = func() martini.Handler {
	return func(token oauth2.Tokens, w http.ResponseWriter, r *http.Request, c martini.Context) {
		if token == nil || token.Expired() {
			next := url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, oauth2.PathLogin+"?next="+next, 302)
			return
		}
		id, err := GetId(token.Access())
		if err != nil {
			log.Printf("Error getting captain token id:", err.Error())
			http.Redirect(w, r, "/error", 302)
			return
		}

		user := ols.GetUserDAO().GetUserFB(id)
		if user.LeagueId == 0 {
			next := url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, "/register?next="+next, 302)
			return
		}
		// TODO - fix this
		if true {
			c.Map(user)
			c.Next()
		} else {
			http.Redirect(w, r, "/captain", 401)
			return
		}

	}
}()

var PlayerRequiredFunc = func() martini.Handler {
	return func(token oauth2.Tokens, w http.ResponseWriter, r *http.Request, c martini.Context) {
		if token == nil || token.Expired() {
			next := url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, oauth2.PathLogin+"?next="+next, 302)
			return
		}
		id, err := GetId(token.Access())
		if err != nil {
			log.Printf("Error getting player token id:", err.Error())
			http.Redirect(w, r, "/error", 302)
			return
		}
		user := ols.GetUserDAO().GetUserFB(id)
		if user.LeagueId == 0 {
			next := url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, "/register?next="+next, 302)
		}

		c.Map(user)
		c.Next()

	}
}()

var DebugPlayerRequired = func() martini.Handler {
	return func(urls url.Values, c martini.Context, w http.ResponseWriter, r *http.Request) {
		leagueIdStr := urls.Get("debug")
		//leagueId, err := strconv.ParseInt(leagueIdStr, 10, 64)

		user := ols.GetUserDAO().GetUserByIgn(leagueIdStr)
		c.Map(user)
	}
}()

var DebugLoginRequired = func() martini.Handler {
	return func(urls url.Values, c martini.Context, w http.ResponseWriter, r *http.Request) {
		loginId := urls.Get("login")
		_, err := GetId(loginId) // Make sure you logged in correctly dope.
		if err != nil {
			panic(err)
		}
		c.MapTo(&DebugToken{urls.Get("login")}, (*oauth2.Tokens)(nil))
	}
}()

// In order to properly emulate Debugging of "Logged in Facebook users" I had to create my own tokens since the framework I was using didn't expose creation X_X
type DebugToken struct {
	Id string
}

// Access returns the access token.
func (t *DebugToken) Access() string {
	return t.Id
}

// Refresh returns the refresh token.
func (t *DebugToken) Refresh() string {
	return ""
}

func (t *DebugToken) Valid() bool {
	return true
}

// Expired returns whether the access token is expired or not.
func (t *DebugToken) Expired() bool {
	return false
}

// ExpiryTime returns the expiry time of the user's access token.
func (t *DebugToken) ExpiryTime() time.Time {
	return time.Now().Add(time.Duration(100) * time.Second)
}

// String returns the string representation of the token.
func (t *DebugToken) String() string {
	return fmt.Sprintf("tokens: %v", t)
}

var CaptainRequired = func() martini.Handler {
	if Debug {
		return DebugPlayerRequired
	} else {
		return CaptainRequiredFunc
	}
}()

var PlayerRequired = func() martini.Handler {
	if Debug {
		return DebugPlayerRequired
	} else {
		return PlayerRequiredFunc
	}
}()

var LoginRequired = func() martini.Handler {
	if Debug {
		return DebugLoginRequired
	} else {
		return oauth2.LoginRequired
	}
}()
