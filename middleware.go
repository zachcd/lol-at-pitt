package main

import (
	"fmt"
	"github.com/go-martini/martini"
	dao "github.com/lab-d8/lol-at-pitt/db"
	"github.com/lab-d8/lol-at-pitt/draft"
	"github.com/lab-d8/oauth2"
	"labix.org/v2/mgo"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// PARAMS is a middleware binder for injecting the params into each handler
func PARAMS(req *http.Request, c martini.Context) {
	req.ParseForm()
	response := req.Form
	c.Map(response)
	c.Next()
}

// DB is a middleware binder that injects the mongo db into each handler
func DB() martini.Handler {
	session, err := mgo.Dial(dao.MongoLocation)
	if err != nil {
		panic(err)
	}

	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB(dao.DatabaseName))
		defer s.Close()
		c.Next()
	}
}

func DRAFT() martini.Handler {
	session, err := mgo.Dial(dao.MongoLocation)
	if err != nil {
		panic(err)
	}
	db := session.DB(dao.DatabaseName)
	olsDraft := draft.Load(db)
	olsDraft.Resume()
	olsDraft.Current.Id = 2

	return func(c martini.Context) {
		c.Map(olsDraft)
		c.Next()
	}
}

var CaptainRequired = func() martini.Handler {
	if Debug {
		return DebugPlayerRequired
	} else {
		return CaptainRequiredFunc
	}
}()

var LoginRequired = func() martini.Handler {
	if Debug {
		return DebugLoginRequired
	} else {
		return oauth2.LoginRequired
	}
}()

var CaptainRequiredFunc = func() martini.Handler {
	return func(token oauth2.Tokens, w http.ResponseWriter, r *http.Request, c martini.Context) {
		if token == nil || token.Expired() {
			next := url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, oauth2.PathLogin+"?next="+next, 302)
		}

		id, err := GetId(token.Access())
		if err != nil {
			http.Redirect(w, r, "/", 304)
		}

		user := dao.GetUserDAO().GetUserFB(id)
		if user.IsCaptain() {
			c.Map(user)
			c.Next()
		} else {
			next := url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, "/register?next="+next, 302)
		}

	}
}()

var DebugPlayerRequired = func() martini.Handler {
	return func(urls url.Values, c martini.Context, w http.ResponseWriter, r *http.Request) {
		leagueIdStr := urls.Get("debug")
		leagueId, err := strconv.ParseInt(leagueIdStr, 10, 64)
		if err != nil {
			panic(err)
		}

		user := dao.GetUserDAO().GetUserLeague(leagueId)
		if user.LeagueId == 0 {
			next := url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, "/register?login="+leagueIdStr+"&next="+next, 302)
		}
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

// In order to preoprly emulate Debugging of "Logged in Facebook users" I had to create my own tokens since the framework I was using didn't expose creation X_X
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
