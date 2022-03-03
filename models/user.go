package models

import (
	"auth/info"
	"context"
	"log"

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
	clientOptions := options.Client().ApplyURI(info.Mongo)
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

//User модель для mongo
type User struct {
	ID      primitive.ObjectID `bson:"_id"`
	GUID    string             `bson:"guid"`
	Refresh string             `bson:"refresh"`
}

type UserInfo struct {
	GUID string `json:"guid"`
}

func (user *User) Create() error {
	_, err := collection.InsertOne(ctx, user)
	return err
}

func (user *User) Get(GUID string) error {
	filter := bson.D{primitive.E{Key: "guid", Value: GUID}, primitive.E{Key: "_id", Value: user.ID}}
	tmp := collection.FindOne(ctx, filter)
	return tmp.Decode(&user)
}

func (user *User) GetByRefrashToken(token string) error {
	filter := bson.D{primitive.E{Key: "refrash", Value: token}}
	tmp := collection.FindOne(ctx, filter)
	return tmp.Decode(&user)
}

func (user *User) UpdateAddNewToken(refrash string) error {
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.D{
			{"$set", bson.D{{"refrash", refrash}}},
		},
	)
	return err
}

// func (user *User) Update(outdatedR, refrash string) error {
// 	filter := bson.D{primitive.E{Key: "refresh", Value: outdatedR}}

// 	update := bson.D{primitive.E{Key: "refrash", Value: refrash}}

// 	return collection.FindOneAndUpdate(ctx, filter, update).Decode(&user)
// }
