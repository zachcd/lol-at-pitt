package main

import (
	"github.com/go-martini/martini"
	"github.com/lab-d8/lol-at-pitt/draft"
	"github.com/lab-d8/lol-at-pitt/ols"
	"github.com/lab-d8/oauth2"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	goauth2 "golang.org/x/oauth2"
	"labix.org/v2/mgo"
	"net/http"
	"net/url"
)

var olsDraft draft.Draft
var ClientId string = "404221166401700"
var ApiSecret string = "bbda73d1673fe517166bf688da56e519"
var Debug bool = true

func main() {
	m := martini.Classic()

	// Setup middleware to be attached to the controllers on every call.
	if Debug {
		InitDebugMiddleware(m)
	} else {
		InitMiddleware(m)
	}

	m.Use(render.Renderer(render.Options{Directory: TemplatesLocation}))

	m.Use(martini.Static("public", martini.StaticOptions{Prefix: "/public"}))

	// TODO: Individual variables not sustainable. Need a better system.
	teamHandler := func(mongo *mgo.Database, renderer render.Render) {
		teams := ols.QueryAllTeams(mongo)
		renderer.HTML(200, "teams", teams)
	}

	individualTeamHandler := func(db *mgo.Database, params martini.Params, renderer render.Render) {
		team := ols.QueryTeam(db, params["name"])
		renderer.HTML(200, "team", team)
	}
	m.Get("/teams", teamHandler)
	m.Get("/team/:name", individualTeamHandler)
	m.Get("/draft", CaptainRequired, func(renderer render.Render, d *draft.Draft) {
		renderer.JSON(200, d)
	})
	m.Get("/draft/bid", func(renderer render.Render, d *draft.Draft) {
		// TODO: Put in CaptainRequired which gets a Player to match the auctioner.
		// TODO: Bid using bid function
	})

	m.Get("/register", oauth2.LoginRequired, func(urls url.Values) {
		urls.Get("name")
		// TODO: Get the summoner name, look it up in Players. If not there, create it.
		// TODO: Create user struct: Id, summoner_id, roles[]
	})
	http.ListenAndServe(":8080", m) // Nginx needs to redirect here, so we don't need sudo priv to test.

}

func InitMiddleware(m *martini.ClassicMartini) {
	m.Use(PARAMS)
	m.Use(DB())
	m.Use(DRAFT())
	m.Use(sessions.Sessions("lol_session", sessions.NewCookieStore([]byte("secret123"))))
	m.Use(oauth2.Facebook(
		&goauth2.Config{
			ClientID:     ClientId,
			ClientSecret: ApiSecret,
			Scopes:       []string{"public_profile", "email", "user_friends"},
			RedirectURL:  "http://local.foo.com/oauth2callback",
		},
	))

}

func InitDebugMiddleware(m *martini.ClassicMartini) {
	m.Use(PARAMS)
	m.Use(DB())
	m.Use(DRAFT())
	m.Use(sessions.Sessions("lol_session", sessions.NewCookieStore([]byte("secret123"))))

}
