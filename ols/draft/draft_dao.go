package draft

import (
	"labix.org/v2/mgo"
)

type DraftDAO struct {
	db         *mgo.Database
	collection *mgo.Collection
}

func (d *DraftDAO) Save(draft *Draft) {
	amt, _ := d.collection.Count()

	if amt == 0 {
		d.collection.Insert(draft)
	} else {
		d.collection.Update(map[string]string{}, draft)
	}
}

func InitDraftDAO(db *mgo.Database) DraftDAO {
	coll := db.C("draft")
	dao := DraftDAO{db, coll}
	return dao
}

func (d *DraftDAO) Load() *Draft {
	var draft Draft
	d.collection.Find(map[string]string{}).One(&draft)
	return &draft
}

func (d *DraftDAO) Delete() {
	d.collection.DropCollection()
}
