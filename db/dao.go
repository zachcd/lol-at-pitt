package db

import (
	"labix.org/v2/mgo"
)

type DAO struct {
	db         *mgo.Database
	Collection *mgo.Collection
}

var (
	usersDAO      *UsersDAO
	matchesDAO    *MatchesDAO
	playersDAO    *PlayersDAO
	DatabaseName  = "lolpitt"
	MongoLocation = "mongodb://localhost"
)

func (d *DAO) Save(needle interface{}, update interface{}) {
	count, _ := d.Collection.Find(needle).Count()
	if count > 0 {
		d.Collection.Update(needle, update)
	} else {
		d.Collection.Insert(update)
	}
}

func GetUserDAO() *UsersDAO {
	if usersDAO == nil || usersDAO.db == nil {
		db := initDB()
		usersDAO = NewUserDAO(db)
	}

	return usersDAO
}

func GetPlayersDAO() *PlayersDAO {
	if playersDAO == nil || playersDAO.db == nil {
		db := initDB()
		playersDAO = NewPlayerContext(db)
	}

	return playersDAO
}

func GetMatchesDAO() *MatchesDAO {
	if matchesDAO == nil || matchesDAO.db == nil {
		db := initDB()
		matchesDAO = NewMatchesContext(db)
	}

	return matchesDAO
}

func initDB() *mgo.Database {
	session, err := mgo.Dial(MongoLocation)
	if err != nil {
		panic(err)
	}
	db := session.DB(DatabaseName)
	return db
}
