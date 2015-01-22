package ols

import (
	"github.com/TrevorSStone/goriot"
	"labix.org/v2/mgo"
)

type PlayersDAO struct {
	DAO
}

func NewPlayerContext(db *mgo.Database) *PlayersDAO {
	dao := PlayersDAO{DAO{db, db.C("players")}}
	return &dao
}

func (p *PlayersDAO) LoadIGN(ign string) Player {
	var player Player
	p.Collection.Find(map[string]string{"ign": ign}).One(&player)
	return player
}

func (p *PlayersDAO) LoadNormalizedIGN(ign string) Player {
	var player Player
	norm := goriot.NormalizeSummonerName(ign)[0]
	p.Collection.Find(map[string]string{"normalizedign": norm}).One(&player)
	return player
}

func (p *PlayersDAO) Load(id int64) Player {
	var player Player
	p.Collection.Find(map[string]int64{"id": id}).One(&player)
	return player
}

func (p *PlayersDAO) All() Players {
	var players Players
	p.Collection.Find(map[string]string{}).All(&players)
	return players
}

func (p *PlayersDAO) Save(player Player) {
	p.DAO.Save(map[string]int64{"id": player.Id}, player)
}

func (p *PlayersDAO) Delete(player Player) {
	p.Collection.Remove(player)
}
