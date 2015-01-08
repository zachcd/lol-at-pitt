package draft

import (
	"github.com/lab-d8/lol-at-pitt/ols"
	"labix.org/v2/mgo"
	"testing"
)

const databaseName string = "lolpitt"
const mongoLocation = "mongodb://localhost"

func TestDraftDAO(t *testing.T) {
	dao := initDAO()
	if dao.collection.Name != "draft" {
		t.Error("Failed, collection name is not draft")
	}

	if dao.db == nil {
		t.Error("Failed, no mongodb")
	}

}

func TestDraftDAOSave(t *testing.T) {
	dao := initDAO()

	draft := Draft{}
	draft.Current = DraftPlayer{1234, 10, "hello", false}
	dao.Save(&draft)

	var savedDraft Draft
	dao.collection.Find(map[string]string{}).One(&savedDraft)
	if savedDraft.Current.Id != draft.Current.Id {
		t.Error("Did not save correctly")
	}

	draft.Current = DraftPlayer{1111, 10, "hello", false}
	dao.Save(&draft)

	dao.collection.Find(map[string]string{}).One(&savedDraft)
	if savedDraft.Current.Id != draft.Current.Id {
		t.Error("Did not update correctly")
	}

	dao.Delete()

}

func TestDraftDAODelete(t *testing.T) {
	dao := initDAO()

	draft := Draft{}
	draft.Current = DraftPlayer{1234, 10, "hello", false}
	dao.Save(&draft)

	dao.Delete()

	savedDraft := dao.Load()
	emptyDraft := Draft{}

	// Annoying, use empty id, since it defaults to empty
	if savedDraft.Current.Id != emptyDraft.Current.Id {
		t.Error("Didn't delete properly", savedDraft)
	}
}
func TestDraftDAOLoad(t *testing.T) {
	dao := initDAO()

	draft := Draft{}
	draft.Current = DraftPlayer{1234, 10, "hello", false}
	dao.Save(&draft)
	savedDraft := dao.Load()
	if savedDraft.Current.Id != draft.Current.Id {
		t.Error("Failed loading properly")
	}

	dao.Delete()
}

func TestDraftInit(t *testing.T) {
	db := initDB()
	playerDAO := ols.NewPlayerContext(db)
	player := ols.Player{}
	player.Id = 1
	player.Ign = "derp"

	playerDAO.Save(player)

	draft := InitNewDraft(db)
	if draft.Current.Id != 1 {
		t.Error("Did not work..")
	}

	playerDAO.Delete(player)
}

func TestDraftBet(t *testing.T) {
	draft := Draft{Current: DraftPlayer{1, -1, "", false}}
	auctioners := map[string]Auctioner{
		"derp": Auctioner{Id: 1, Points: 50, Team: "derp"},
		"lol":  Auctioner{Id: 2, Points: 60, Team: "lol"},
	}

	draft.Auctioners = auctioners

	badBid := draft.Bid(90, "derp")
	badBid2 := draft.Bid(1, "llll")
	goodBid := draft.Bid(10, "derp")
	littleBid := draft.Bid(9, "lol")
	copyBid := draft.Bid(11, "derp")
	goodBid2 := draft.Bid(11, "lol")

	if badBid || littleBid || copyBid || badBid2 {
		t.Error("Bad bids went through")
	}
	if !(goodBid && goodBid2) {
		t.Error("Good bids didn't go through..")
	}

}

func initDAO() DraftDAO {
	db := initDB()
	dao := InitDraftDAO(db)
	return dao
}
func initDB() *mgo.Database {
	session, _ := mgo.Dial(mongoLocation)
	db := session.DB(databaseName)
	return db
}
