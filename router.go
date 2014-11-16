package main

import (
	"github.com/codegangsta/martini"
	goauth2 "github.com/golang/oauth2"
	"github.com/lab-d8/lol-at-pitt/ols"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"labix.org/v2/mgo"
	"net/http"
	"net/url"
)

func main() {
	m := martini.Classic()
	// Setup middleware to be attached to the controllers on every call.
	m.Use(DB())
	m.Use(sessions.Sessions("lol_session", sessions.NewCookieStore([]byte("secret123"))))
	m.Use(render.Renderer(render.Options{Directory: TemplatesLocation}))
	m.Use(PARAMS)
	m.Use(martini.Static("public", martini.StaticOptions{Prefix: "/public"}))
	m.Use(oauth2.Facebook(
		goauth2.Client(FBClientId, FBApiSecret),
		goauth2.RedirectURL("http://local.foo.com/oauth2callback"),
		goauth2.Scope("public_profile"),
	))

	// Test the login functionality!
	m.Get("/success", oauth2.LoginRequired, func() string {
		return "You are logged in!"
	})

	// TODO: Individual variables not sustainable. Need a better system.
	teamHandler := func(mongo *mgo.Database, urls url.Values, renderer render.Render) {
		teams := ols.QueryAllTeams(mongo)
		renderer.HTML(200, "teams", teams)
	}

	individualTeamHandler := func(db *mgo.Database, params martini.Params, renderer render.Render) {
		team := ols.QueryTeam(db, params["name"])
		renderer.HTML(200, "team", team)
	}
	m.Get("/teams", teamHandler)
	m.Get("/team/:name", individualTeamHandler)

	http.ListenAndServe(":8080", m) // Nginx needs to redirect here, so we don't need sudo priv to test.

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
