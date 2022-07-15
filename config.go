package mongoose

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client
var DatabaseName string

func ConnectDB(uri string, databaseName string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))

	DatabaseName = databaseName

	if err != nil {
		log.Fatal("Could not initialize client: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatal("Could not connect to client: ", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Could not ping the database: ", err)
	}

	fmt.Println("Connected to MongoDB")
	DB = client

	return client
}

func GetCollection(collectionName string) *mongo.Collection {
	collection := DB.Database(DatabaseName).Collection(collectionName)
	return collection
}

func CreateIndex(collectionName string, field string, unique bool, sparse bool) bool {

	// 1. Lets define the keys for the index we want to create
	mod := mongo.IndexModel{
		Keys:    bson.M{field: 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(unique).SetSparse(sparse),
	}

	// 2. Create the context for this operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Connect to the database and access the collection
	collection := DB.Database(DatabaseName).Collection(collectionName)

	// 4. Create a single index
	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		// 5. Something went wrong, we log it and return false
		fmt.Println(err.Error())
		return false
	}

	// 6. All went well, we return true
	return true
}
