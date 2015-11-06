package managers
import (
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
)

var dbSession *mgo.Session
var LogLevel int = 1

/*
Log Levels:
0 - debug
1 - info
2 - access
3 - warn
4 - error
 */

const DB_HOST string = "localhost"
const DB_PORT string = "27017"
const DB_NAME string = "CodeCollaborate"

func ConnectMGo() {

	var err error
	dbSession, err = mgo.Dial(DB_HOST + ":" + DB_PORT + "/" + DB_NAME)
	if err != nil {
		log.Fatal("Error connecting to MongoDB instance: ", err)
	}

	dbSession.SetMode(mgo.Strong, true)
	//	dbSession.SetMode(mgo.Eventual, true)
	LogInfo("Connected to DB")
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

func LogError(message string, err error) {
	if (LogLevel <= 4) {

		session, collection := GetMGoCollection("Log")
		defer session.Close()

		collection.Insert(bson.M{
			"level": "error",
			"message": message,
			"error": err.Error(),
			"timestamp": time.Now().UTC().String(),
		})
	}
}

func LogWarn(message string) {
	if (LogLevel <= 3) {

		session, collection := GetMGoCollection("Log")
		defer session.Close()

		collection.Insert(bson.M{
			"level": "warn",
			"message": message,
			"timestamp": time.Now().UTC().String(),
		})

		log.Println(message)
	}
}

func LogAccess(baseRequest baseRequests.BaseRequest, rawMessage string) {
	if (LogLevel <= 2) {
		session, collection := GetMGoCollection("Log")
		defer session.Close()

		collection.Insert(bson.M{
			"level": "access",
			"resource": baseRequest.Resource,
			"resource_id": baseRequest.ResId,
			"username": baseRequest.Username,
			"timestamp": time.Now().UTC().String(),
			"raw_message": rawMessage,
		})
	}
}

func LogInfo(message string) {
	if (LogLevel <= 1) {

		session, collection := GetMGoCollection("Log")
		defer session.Close()

		collection.Insert(bson.M{
			"level": "info",
			"message": message,
			"timestamp": time.Now().UTC().String(),
		})
	}
}

func LogDebug(message string) {
	if (LogLevel <= 0) {

		session, collection := GetMGoCollection("Log")
		defer session.Close()

		collection.Insert(bson.M{
			"level": "debug",
			"message": message,
			"timestamp": time.Now().UTC().String(),
		})
	}
}

func NewObjectIdString() string {
	return bson.NewObjectId().Hex()
}
