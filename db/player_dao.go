package db

import (
	"github.com/lab-d8/lol-at-pitt/ols"
	"labix.org/v2/mgo"
)

type PlayersDAO struct {
	DAO
}

func NewPlayerContext(db *mgo.Database) *PlayersDAO {
	dao := PlayersDAO{DAO{db, db.C("players")}}
	return &dao
}

func (p *PlayersDAO) LoadIGN(ign string) ols.Player {
	var player ols.Player
	p.Collection.Find(map[string]string{"ign": ign}).One(&player)
	return player
}

func (p *PlayersDAO) Load(id int64) ols.Player {
	var player ols.Player
	p.Collection.Find(map[string]int64{"id": id}).One(&player)
	return player
}

func (p *PlayersDAO) All() ols.Players {
	var players ols.Players
	p.Collection.Find(map[string]string{}).All(&players)
	return players
}

func (p *PlayersDAO) Save(player ols.Player) {
	p.DAO.Save(map[string]int64{"id": player.Id}, player)
}

func (p *PlayersDAO) Delete(player ols.Player) {
	p.Collection.Remove(player)
}
