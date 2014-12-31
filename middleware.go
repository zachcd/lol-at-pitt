package main

import (
	"github.com/go-martini/martini"
	"github.com/lab-d8/lol-at-pitt/draft"
	"github.com/lab-d8/oauth2"
	"labix.org/v2/mgo"
	"net/http"
	"net/url"
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
	session, err := mgo.Dial(MongoLocation)
	if err != nil {
		panic(err)
	}

	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB(DatabaseName))
		defer s.Close()
		c.Next()
	}
}

func DRAFT() martini.Handler {
	session, err := mgo.Dial(MongoLocation)
	if err != nil {
		panic(err)
	}
	db := session.DB(DatabaseName)
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

var CaptainRequiredFunc = func() martini.Handler {
	return func(token oauth2.Tokens, w http.ResponseWriter, r *http.Request) {
		if token == nil || token.Expired() {
			next := url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, oauth2.PathLogin+"?next="+next, 302)
		}
	}
}()

var DebugPlayerRequired = func() martini.Handler {
	return func(urls url.Values) {
		name := urls.Get("debug")
		_ = name
	}
}()
