package userModels

import (
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

)

type User struct {
	_id           string // ID of object
	Username      string // Username
	Email         string // Email of user
}



func GetUser(session *mgo.Session, id string) (*User, error) {
	copySession := session.Copy()
	defer copySession.Close()

	collection := copySession.DB("").C("Files")

	result := new(User)
	err := collection.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		log.Println("Failed to retrieve User")
		log.Println(err)
		return nil, err
	}

	return result, nil
}