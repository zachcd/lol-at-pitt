package main

import (
	"github.com/TrevorSStone/goriot"
	"github.com/go-martini/martini"
	dao "github.com/lab-d8/lol-at-pitt/db"
	"github.com/lab-d8/lol-at-pitt/draft"
	"github.com/lab-d8/lol-at-pitt/ols"
	"github.com/lab-d8/lol-at-pitt/site"
	"github.com/lab-d8/oauth2"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	goauth2 "golang.org/x/oauth2"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"net/url"
	"sort"
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

	teamHandler := func(mongo *mgo.Database, renderer render.Render) {
		teams := ols.QueryAllTeams(mongo)
		renderer.HTML(200, "teams", teams)
	}

	individualTeamHandler := func(db *mgo.Database, params martini.Params, renderer render.Render) {
		team := ols.QueryTeam(db, params["name"])
		renderer.HTML(200, "team", team)
	}

	m.Get("/error", func(urls url.Values, renderer render.Render) {
		renderer.HTML(200, "error", urls.Get("status"))
	})
	m.Get("/teams", teamHandler)
	m.Get("/team/:name", individualTeamHandler)

	m.Get("/captain", CaptainRequired, func(user site.User, renderer render.Render) {
		renderer.JSON(200, user)
	})

	m.Get("/player", PlayerRequired, func(renderer render.Render, user site.User) {
		player := dao.GetPlayersDAO().Load(user.LeagueId)
		renderer.JSON(200, player)

	})

	m.Get("/register", LoginRequired, func(urls url.Values, renderer render.Render) {
		renderer.HTML(200, "register", Register{Next: urls.Get("next")})
	})

	m.Get("/oauth2error", func(token oauth2.Tokens, renderer render.Render) {
		renderer.JSON(200, token)
	})

	m.Get("/draft/summary", func(renderer render.Render) {
		allPlayers := dao.GetPlayersDAO().All()
		players := allPlayers.Filter(func(player ols.Player) bool {
			return !player.Captain && player.Team != ""
		})

		sort.Sort(players)
		renderer.HTML(200, "drafted", players)
	})

	m.Get("/register/complete", LoginRequired, func(urls url.Values, renderer render.Render, token oauth2.Tokens, w http.ResponseWriter, r *http.Request) {
		summonerName := urls.Get("summoner")
		normalizedSummonerName := goriot.NormalizeSummonerName(summonerName)[0]
		result, err := goriot.SummonerByName("na", normalizedSummonerName)

		if err != nil || token.Expired() {
			http.Redirect(w, r, "/error?status=InvalidFacebook", 302)
			return
		}

		id, err := GetId(token.Access())
		if err != nil {
			renderer.Status(404)
			return
		}
		//Check for captain
		summonerProfile := result[normalizedSummonerName]
		player := dao.GetPlayersDAO().Load(summonerProfile.ID)
		user := dao.GetUserDAO().GetUserFB(id)

		// User is registered registered
		if user.LeagueId != 0 {
			http.Redirect(w, r, "/error?status=AlreadyRegistered", 302)
			return
		}
		user = site.User{LeagueId: summonerProfile.ID, FacebookId: id}
		log.Println("User registered:", user)
		if player.Id == 0 {
			// new player not in our db
			player := ols.Player{Id: summonerProfile.ID, Ign: summonerProfile.Name}
			log.Println("New player:", player)
			dao.GetPlayersDAO().Save(player)
		}
		dao.GetUserDAO().Save(user)
		//next := urls.Get("next")
		log.Println("register completed going to page?")
		renderer.HTML(200, "register_complete", 1)
	})

	err := http.ListenAndServe(":6060", m) // Nginx needs to redirect here, so we don't need sudo priv to test.
	if err != nil {
		log.Println(err)
	}

}

func InitMiddleware(m *martini.ClassicMartini) {
	m.Use(PARAMS)
	m.Use(DB())
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
	m.Use(sessions.Sessions("lol_session", sessions.NewCookieStore([]byte("secret123"))))
	m.Use(render.Renderer(render.Options{Directory: TemplatesLocation}))
	m.Use(martini.Static("public", martini.StaticOptions{Prefix: "/public"}))
	SetId("1", "10153410152015744") // Me
}
