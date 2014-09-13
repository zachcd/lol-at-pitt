package db

import (
	"labix.org/v2/mgo"
)

// Insert takes a name and an object and will put it in mongo under thus.
func Insert(db *mgo.Database, collectionName string, inserted interface{}) {
	err := db.C(collectionName).Insert(inserted)

	if err != nil {
		panic("Db error")
	}
}

// Update uses the objectMatch to find the old object and objectUpdated to update with new data.
// Generally these are the same objects with just updated information. Unique keys that identify them should be the same
func Update(db *mgo.Database, collectionName string, objectMatch interface{}, objectUpdated interface{}) {
	err := db.C(collectionName).Update(objectMatch, objectUpdated)
	if err != nil {
		panic("Db error")
	}
}

// Remove takes out an object from a collection with the given information.
// Again, it is the responsibility of the programmer using this to make sure that obj has a unique identifier for mgo to remove from
func Remove(db *mgo.Database, collectionName string, obj interface{}) {
	err := db.C(collectionName).Remove(obj)

	if err != nil {
		panic("Db error")
	}
}

//Query Fill the entire object passed in
func Query(db *mgo.Database, collectionName string, obj interface{}) {
	db.C(collectionName).Find(obj).One(obj)
}

//Count just counts.
func Count(db *mgo.Database, collectionName string, obj interface{}) int {
	count, _ := db.C(collectionName).Find(obj).Count()
	return count
}
