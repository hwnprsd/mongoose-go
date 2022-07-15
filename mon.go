package mongoose

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	Users                  = "users"
	Payments               = "payments"
	QuizTemplates          = "quiz_templates"
	QuizEntries            = "quiz_entries"
	Campaigns              = "campaigns"
	CampaignParticipations = "campaign_participations"
	Rewards                = "rewards"
)

type Models interface {
}

// Create a safe context with no timeouts
// TODO: Extend this func to accept contexts from the calling functions
func GetContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx, cancel
}

// Create an ObjectID from a string
func ObjId(s string) primitive.ObjectID {
	o, _ := primitive.ObjectIDFromHex(s)
	return o
}

// Create a new DataTime Mongo Object from the current date
func Now() primitive.DateTime {
	return primitive.NewDateTimeFromTime(time.Now())
}

// Struct to store the meta data of the struct
// Stores the raw mongo instance of the model, as well as the model name given during instantiation
// Used bind functions to the `A` type, so the pointers have valid types
type CollectionWrapper[A Models] struct {
	I    mongo.Collection // Instance of the collection
	name string
}

// Constructor to create a new Model, given the type of model and the collection name in string
// Returns a reference to the Collection Wrapper, which contains helpful methods
// Initialize this when you inititalize your models
// Example
// userWrapper := mongo.New[User]('users')
func NewCollectionWrapper[A Models](name string) *CollectionWrapper[A] {
	return &CollectionWrapper[A]{
		I:    *GetCollection(name),
		name: name,
	}
}

// FindOne query insprired by mongoose in idomatic go
//
// Example
//
// 		user := User{}
// 		if err := userWrapper.FindOne(bson.M{"Name": "king"}, &user); err != nil {
//			panic(err)
//		}
func (c CollectionWrapper[M]) FindOne(query bson.M, model *M) error {
	ctx, cancel := GetContext()
	defer cancel()
	if err := c.I.FindOne(ctx, query).Decode(model); err != nil {
		return errors.New("Cannot find document")
	}
	return nil
}

// FindOneById query insprired by mongoose in idomatic go
//
// Example
//
// 		user := User{}
// 		if err := userWrapper.FindOneById("USER_ID", &user); err != nil {
//			panic(err)
//		}
func (c CollectionWrapper[M]) FindOneById(id string, model *M) error {
	return c.FindOne(bson.M{"_id": ObjId(id)}, model)
}

// FindByIdAndUpdate query insprired by mongoose in idomatic go
//
// Example
//
//			updatedUser := User{}
//			update := bson.M{
//			"$set": bson.M{
//					"lastName": "Queen",
//				},
//			}
// 			if err := userWrapper.FindOneById("USER_ID", update, &updatedUser); err != nil {
//				panic(err)
//			}
func (c CollectionWrapper[M]) FindByIdAndUpdate(idHex string, update bson.M, updated *M) error {
	return c.FindOneAndUpdate(bson.M{"_id": ObjId(idHex)}, update, updated)
}

// FindByIdAndUpdate query insprired by mongoose in idomatic go
//
// Example
//
//			updatedUser := User{}
//			update := bson.M{
//			"$set": bson.M{
//					"Name": "Queen",
//				},
//			}
// 			if err := userWrapper.FindOne(bson.M{"Name": "king"}, update, &updatedUser); err != nil {
//				panic(err)
//			}
func (c CollectionWrapper[M]) FindOneAndUpdate(find bson.M, update bson.M, updated *M) error {
	ctx, cancel := GetContext()
	defer cancel()
	if err := c.I.FindOneAndUpdate(ctx, find, update).Decode(updated); err != nil {
		fmt.Println(err)
		return errors.New("Cannot update document")
	}
	return nil
}

// A clean wrapper around creating a new document
//
// Example
//
// 	user := User{
// 		"Name": "King",
// 		"Role": "King"
// 	}
// 	if err := userWrapper.New(user); err != nil {
// 		panic(err)
// 	}
func (c CollectionWrapper[M]) New(document M) error {
	ctx, cancel := GetContext()
	defer cancel()
	_, err := c.I.InsertOne(ctx, document)
	return err
}

// FindMany inspired by Mongoose. Find a slice of documents instead of just one
//
// Example
//
// 	users := []User{} // Slice of users
// 	query := bson.M{
// 		"Role": "King"
// 	}
// 	if err := userWrapper.FindMany(query, &users); err != nil {
// 		panic(err)
// 	}
func (c CollectionWrapper[M]) FindMany(bson interface{}, elem *[]M) error {
	ctx, cancel := GetContext()
	defer cancel()
	cursor, err := c.I.Find(ctx, bson)
	if err != nil {
		return err
	}
	if err := cursor.All(ctx, elem); err != nil {
		return err
	}
	return nil
}

// Struct to help populate fields
type Populate struct {
	LocalField, ForeignModel, As string
}

// Find many and populate using mongo driver's aggregate
//
// Example
//
// 	query := bson.D{{
// 		Key:   "_id",
// 		Value: models.ObjId(campaignId),
// 	}}
//
// 	populate := mongo.Populate{
// 		As:           "Quizzes.Data",
// 		ForeignModel: models.QuizTemplates,
// 		LocalField:   "Quizzes.Ids",
// 	}
//
// 	campaigns := []Campaign{}
// 	if err := campaignWrapper.FindManyPopulate(query, populate, &campaigns); err != nil {
// 		utils.Panic(401, "[1] Campaign Not Found", err)
// 	}
func (c CollectionWrapper[M]) FindManyPopulate(matchQuery bson.D, populate Populate, elem *[]M) error {
	ctx, cancel := GetContext()
	defer cancel()

	match := bson.D{{Key: "$match", Value: matchQuery}}
	lookup := bson.D{{Key: "$lookup", Value: bson.D{{
		Key:   "from",
		Value: populate.ForeignModel,
	}, {
		Key:   "localField",
		Value: populate.LocalField,
	}, {
		Key:   "foreignField",
		Value: "_id",
	}, {
		Key:   "as",
		Value: populate.As,
	}}}}

	cursor, err := c.I.Aggregate(ctx, mongo.Pipeline{match, lookup})
	if err != nil {
		return err
	}
	if err := cursor.All(ctx, elem); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
