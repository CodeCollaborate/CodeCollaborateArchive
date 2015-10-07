package userModels

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"gopkg.in/mgo.v2"
	"github.com/CodeCollaborate/CodeCollaborate/modules/user/requests"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
	"strings"
)

type UserAuthData struct {
	Username      string   // Username
	Email         string   // Email of user
	Password_Hash string   // Hashed Password
	Tokens        []string // Token after logged in.
}

func Register(session *mgo.Session, registrationRequest userRequests.UserRegisterRequest) error {

	// Hash password using bcrypt
	pwHashBytes, err := bcrypt.GenerateFromPassword([]byte(registrationRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash password")
		log.Println(err)
		return err;
	}

	// Create new UserAuthData object
	userAuthData := new(UserAuthData)
	userAuthData.Username = registrationRequest.Username
	userAuthData.Email = registrationRequest.Email
	userAuthData.Password_Hash = string(pwHashBytes[:])

	// Get new DB connection
	copySession := session.Copy()
	defer copySession.Close()

	// Get collection
	collection := copySession.DB("").C("Users")

	// Make sure index is unique
	index := mgo.Index{
		Key: []string{"username"},
		Unique:true,
		DropDups:true,
		Background:true,
		Sparse:true,
	}
	err = collection.EnsureIndex(index);
	if err != nil {
		log.Println("Failed to create index in Users collection")
		log.Println(err)
		return err
	}

	// Register new user
	err = collection.Insert(userAuthData)
	if err != nil {
		if !mgo.IsDup(err) {
			log.Println("Failed to register User entry")
			log.Println(err)
		}
		return err
	}

	return nil
}

func Login(session *mgo.Session, loginRequest userRequests.UserLoginRequest) (string, error) {

	// Get new DB connection
	copySession := session.Copy()
	defer copySession.Close()

	// Get collection
	collection := copySession.DB("").C("Users")

	userAuthData := UserAuthData{}
	if err := collection.Find(bson.M{"$or": []interface{}{bson.M{"username": loginRequest.UsernameOREmail}, bson.M{"email": loginRequest.UsernameOREmail}}}).One(&userAuthData); err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userAuthData.Password_Hash), []byte(loginRequest.Password)); err != nil {
		return "", err
	}

	tokenBytes, err := bcrypt.GenerateFromPassword([]byte(loginRequest.UsernameOREmail + time.Now().String()), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to generate token")
		log.Println(err)
		return "", err;
	}

	token := string(tokenBytes[:])

	err = addToken(collection, userAuthData, token)
	if err != nil {
		log.Println("Failed to save token")
		log.Println(err)
		return "", err;
	}

	return token, nil
}

func CheckAuth(session *mgo.Session, baseRequest base.BaseRequest) bool {

	// Get new DB connection
	copySession := session.Copy()
	defer copySession.Close()

	// Get collection
	collection := copySession.DB("").C("Users")

	userAuthData := UserAuthData{}
	if err := collection.Find(bson.M{"$or": []interface{}{bson.M{"username": baseRequest.Username}, bson.M{"email": baseRequest.Username}}}).One(&userAuthData); err != nil {
		return false
	}

	for _, token := range userAuthData.Tokens {
		if (strings.Compare(token, baseRequest.Token) == 0) {
			return true
		}
	}

	return false
}

func addToken(collection *mgo.Collection, userAuthData UserAuthData, token string) error {
	userAuthData.Tokens = append(userAuthData.Tokens, token)

	return collection.Update(bson.M{"username": userAuthData.Username}, bson.M{"$set": bson.M{"tokens": userAuthData.Tokens}})
}