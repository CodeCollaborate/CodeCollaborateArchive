package userModels

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"gopkg.in/mgo.v2"
	"github.com/CodeCollaborate/CodeCollaborate/modules/user/requests"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
	"github.com/CodeCollaborate/CodeCollaborate/managers"
)

type UserAuthData struct {
	Username      string   // Username
	Email         string   // Email of user
	Password_Hash string   // Hashed Password
	Tokens        []string // Token after logged in.
}

func Register(registrationRequest userRequests.UserRegisterRequest) base.WSResponse {

	// Hash password using bcrypt
	pwHashBytes, err := bcrypt.GenerateFromPassword([]byte(registrationRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash password")

		return base.NewFailResponse(-102, registrationRequest.BaseMessage.Tag, nil)
	}

	// Create new UserAuthData object
	userAuthData := new(UserAuthData)
	userAuthData.Username = registrationRequest.Username
	userAuthData.Email = registrationRequest.Email
	userAuthData.Password_Hash = string(pwHashBytes[:])

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Users")
	defer session.Close()

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
		log.Println(err)
		return base.NewFailResponse(-102, registrationRequest.BaseMessage.Tag, nil)
	}

	// Register new user
	err = collection.Insert(userAuthData)
	if err != nil {
		if !mgo.IsDup(err) {
			log.Println(err)
			return base.NewFailResponse(-103, registrationRequest.BaseMessage.Tag, nil)
		}
		return base.NewFailResponse(-102, registrationRequest.BaseMessage.Tag, nil)
	}

	return base.NewSuccessResponse(registrationRequest.BaseMessage.Tag, nil)
}

func Login(loginRequest userRequests.UserLoginRequest) base.WSResponse {

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Users")
	defer session.Close()

	userAuthData := UserAuthData{}
	if err := collection.Find(bson.M{"$or": []interface{}{bson.M{"username": loginRequest.UsernameOREmail}, bson.M{"email": loginRequest.UsernameOREmail}}}).One(&userAuthData); err != nil {
		return base.NewFailResponse(-105, loginRequest.BaseMessage.Tag, nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userAuthData.Password_Hash), []byte(loginRequest.Password)); err != nil {
		return base.NewFailResponse(-105, loginRequest.BaseMessage.Tag, nil)
	}

	tokenBytes, err := bcrypt.GenerateFromPassword([]byte(loginRequest.UsernameOREmail + time.Now().String()), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to generate token")
		log.Println(err)
		return base.NewFailResponse(-104, loginRequest.BaseMessage.Tag, nil)
	}

	token := string(tokenBytes[:])

	err = addToken(collection, userAuthData, token)
	if err != nil {
		log.Println("Failed to save token")
		log.Println(err)
		return base.NewFailResponse(-104, loginRequest.BaseMessage.Tag, nil)
	}

	return base.NewSuccessResponse(loginRequest.BaseMessage.Tag, map[string]interface{}{"token":token})
}

func CheckAuth(baseRequest base.BaseRequest) bool {

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Users")
	defer session.Close()

	userAuthData := UserAuthData{}
	if err := collection.Find(bson.M{"username": baseRequest.Username}).One(&userAuthData); err != nil {
		return false
	}

	for _, token := range userAuthData.Tokens {
		if (token == baseRequest.Token) {
			return true
		}
	}

	return false
}

func addToken(collection *mgo.Collection, userAuthData UserAuthData, token string) error {
	userAuthData.Tokens = append(userAuthData.Tokens, token)

	return collection.Update(bson.M{"username": userAuthData.Username}, bson.M{"$set": bson.M{"tokens": userAuthData.Tokens}})
}