package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var fuck int = 10757

func initFunnyRouter(m *martini.ClassicMartini) {
	m.Get("/fuck/smegs", func(renderer render.Render) {
		renderer.HTML(200, "fuck", fuck)
	})

	m.Get("/fuck/smeg/count", func(renderer render.Render) {
		fuck += 1
		renderer.Redirect("/fuck/smegs", 302)
	})
}
