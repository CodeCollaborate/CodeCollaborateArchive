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


//func main() {
//	session, err := mgo.Dial("localhost/CodeCollaborate")
//	if err != nil {
//		panic(err)
//	}
//	defer session.Close()
//
//	// Optional. Switch the session to a monotonic behavior.
//
//	session.SetMode(mgo.Monotonic, true)
//	c := session.DB("").C("people")
//	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
//		&Person{"Cla", "+55 53 8402 8510"})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	result := Person{}
//	err = c.Find(bson.M{"name": "Ale"}).One(&result)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println("Phone:", result.Phone)
//}