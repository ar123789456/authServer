package models

import (
	"context"
	"log"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.TODO()

const (
	user = "user"
	db   = "auth"
)

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database(db).Collection(user)
}

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Password string             `bson:"password"`
	Refresh  string             `bson:"refresh"`
}

type UserInfo struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func (user *User) Create() error {
	_, err := collection.InsertOne(ctx, user)
	return err
}

func (user *User) Get(name string) error {
	filter := bson.D{primitive.E{Key: "name", Value: name}}
	tmp := collection.FindOne(ctx, filter)
	return tmp.Decode(&user)
}

func (user *User) UpdateAddNewToken(refrash string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: user.ID}}

	update := bson.D{primitive.E{Key: "refrash", Value: refrash}}

	return collection.FindOneAndUpdate(ctx, filter, update).Decode(&user)
}

func (user *User) Update(outdatedR, refrash string) error {
	filter := bson.D{primitive.E{Key: "refresh", Value: outdatedR}}

	update := bson.D{primitive.E{Key: "refrash", Value: refrash}}

	return collection.FindOneAndUpdate(ctx, filter, update).Decode(&user)
}
