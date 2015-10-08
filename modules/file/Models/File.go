package file

import (
	"github.com/CodeCollaborate/CodeCollaborate/modules/project/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type File struct {
	_id           string // ID of object
	Name          string // Name of file
	Relative_Path string // Path of file
	Version       int    // File version
	Project       project.Project // Reference to Project object
}

func (file File) Save(session *mgo.Session){
	copySession := session.Copy()
	defer copySession.Close()

	collection := copySession.DB("").C("Files")
	err := collection.Insert(file)
	if err != nil {
		log.Println(err)
	}
}

func GetFile (session *mgo.Session, id string){
	copySession := session.Copy()
	defer copySession.Close()

	collection := copySession.DB("").C("Files")

	result := File{}
	err := collection.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		log.Println(err)
	}
}