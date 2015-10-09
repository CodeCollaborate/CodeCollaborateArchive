package managers
import (
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)
var dbSession *mgo.Session

const DB_HOST string = "localhost"
const DB_PORT string = "27017"
const DB_NAME string = "CodeCollaborate"

func ConnectMGo() error {

	var err error
	dbSession, err = mgo.Dial(DB_HOST + ":" + DB_PORT + "/" + DB_NAME)
	if err != nil {
		return err
	}

	dbSession.SetMode(mgo.Strong, true)
//	dbSession.SetMode(mgo.Eventual, true)
	log.Println("Connected to DB")

	return nil
}

func GetPrimaryMGoSession() *mgo.Session {
	return dbSession
}

/**
	REMEMBER TO CLOSE THIS RESOURCE CONNECTION WHEN FINISHED.
 */
func GetNewMGoSession() *mgo.Session {
	return dbSession.Copy()
}

/**
	REMEMBER TO CLOSE THIS RESOURCE CONNECTION WHEN FINISHED.
 */
func GetMGoDatabase(dbName string) (*mgo.Session, *mgo.Database) {
	copySession := GetNewMGoSession()

	// Get collection
	return copySession, copySession.DB(dbName)
}

/**
	REMEMBER TO CLOSE THIS RESOURCE CONNECTION WHEN FINISHED.
 */
func GetMGoCollection(collectionName string) (*mgo.Session, *mgo.Collection) {
	// Get collection
	session, database := GetMGoDatabase("")
	return session, database.C(collectionName)
}

func LogError(err error){
	session, collection := GetMGoCollection("errors")
	defer session.Close()

	collection.Insert(bson.M{"error": err.Error(), "timestamp": time.Now().UTC().String()})
}

func NewObjectIdString() string{
	return bson.NewObjectId().Hex()
}