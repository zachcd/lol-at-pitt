package main

import (
	"github.com/TrevorSStone/goriot"
	"github.com/go-martini/martini"
	"github.com/lab-d8/lol-at-pitt/draft"
	"github.com/lab-d8/lol-at-pitt/ols"
	"github.com/lab-d8/lol-at-pitt/site"
	"github.com/lab-d8/oauth2"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	goauth2 "golang.org/x/oauth2"
	"labix.org/v2/mgo"
	"net/http"
	"net/url"
	"time"
)

var olsDraft draft.Draft

type Register struct {
	Id   string
	Name string
	Next string
}

func main() {
	m := martini.Classic()
	goriot.SetAPIKey(LeagueApiKey)
	goriot.SetLongRateLimit(LongLeagueLimit, 10*time.Minute)
	goriot.SetSmallRateLimit(ShortLeagueLimit, 10*time.Second)

	// Setup middleware to be attached to the controllers on every call.
	if Debug {
		InitDebugMiddleware(m)
	} else {
		InitMiddleware(m)
	}

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

	m.Get("/captain", CaptainRequired, func(user site.User, renderer render.Render) {
		renderer.JSON(200, user)
	})
	m.Get("/draft/bid", func(renderer render.Render, d *draft.Draft) {
		// TODO: Put in CaptainRequired which gets a Player to match the auctioner.
		// TODO: Bid using bid function
	})

	m.Get("/register", LoginRequired, func(urls url.Values, renderer render.Render) {
		renderer.HTML(200, "register", Register{Next: urls.Get("next")})
	})

	m.Get("/oauth2error", func(token oauth2.Tokens, renderer render.Render) {
		renderer.JSON(200, token)
	})

	m.Get("/register/create", LoginRequired, func(urls url.Values, renderer render.Render, token oauth2.Tokens) {
		summonerName := urls.Get("summoner")
		normalizedSummonerName := goriot.NormalizeSummonerName(summonerName)[0]
		result, err := goriot.SummonerByName("na", normalizedSummonerName)

		if err != nil {
			renderer.Status(404)
		}
		summonerProfile := result[normalizedSummonerName]
		player := ols.Player{}
		player.Id = summonerProfile.ID

	})

	m.Get("/register/captain")

	http.ListenAndServe(":6060", m) // Nginx needs to redirect here, so we don't need sudo priv to test.

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
			RedirectURL:  "http://www.lol-at-pitt.com/oauth2callback",
		},
	))
	m.Use(render.Renderer(render.Options{Directory: TemplatesLocation}))
	m.Use(martini.Static("public", martini.StaticOptions{Prefix: "/public"}))
}

func InitDebugMiddleware(m *martini.ClassicMartini) {
	m.Use(PARAMS)
	m.Use(DB())
	m.Use(DRAFT())
	m.Use(sessions.Sessions("lol_session", sessions.NewCookieStore([]byte("secret123"))))
	m.Use(render.Renderer(render.Options{Directory: TemplatesLocation}))
	m.Use(martini.Static("public", martini.StaticOptions{Prefix: "/public"}))
	SetId("1", "10153410152015744") // Me
}
