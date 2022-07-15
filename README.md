## MONGOOSE

MongooseGo is a mongodb driver wrapper, which provides a lot of convenience over existing mongodb driver functions.
This package is inspired by mongoose and strives to achieve the same featureset as mongoose from the JS/TS realm

### Installation

Built for Go 1.18 and above

```sh
go mod get github.com/d3fkon/mongoose-go
```

### Usage

#### Connect to a mongodb cluster

```go
mongoose.ConnectDB("mongodb://localhost/my_database", "my_database")
```

#### Defining a model

Define a model using the same way you would for the official mongodb-driver

```go
type User struct {
  FirstName      string `bson:"first_name"`
  LastName       string `bson:"last_name"`
  PassportNumber string `bson:"passport_number"
}
```

Create a mongoose collection wrapper for this model

```go
collectionName := "users"
userWrapper := mongoose.NewCollectionWrapper[User](collectionName)
```

Now through this wrapper, you will have access to the raw instance, on which you can call the native mongodb-driver functions, as well as have the convenience functions on top

```go
userWrapper.I // Accessing the raw mongodb instance
```

Indexing your collection with relevant indexes

```go
mongoose.CreateIndex(collectionName, "passport_number", true, false)
//                   Collection Name, Field Name      ,Unique, Sparse
```

Access different methods

Find One

```go
user := User{}
if err := userWrapper.FindOne(bson.M{"Name": "king"}, &user); err != nil {
  panic(err)
}
```

Find One By Id

```go
user := User{}
if err := userWrapper.FindOneById("USER_ID", &user); err != nil {
  panic(err)
}

```

Find Many

```go
users := []User{} // Slice of users
query := bson.M{
  "Role": "King"
}
if err := userWrapper.FindMany(query, &users); err != nil {
  panic(err)
}
```

Find By Id And Update

```go
  updatedUser := User{}
  update := bson.M{
  "$set": bson.M{
      "lastName": "Queen",
    },
  }
  if err := userWrapper.FindOneById("USER_ID", update, &updatedUser); err != nil {
    panic(err)
  }
```

Find One and Update

```go
  updatedUser := User{}
  update := bson.M{
  "$set": bson.M{
      "Name": "Queen",
    },
  }
  if err := userWrapper.FindOne(bson.M{"Name": "king"}, update, &updatedUser); err != nil {
    panic(err)
  }
```

New

```go

user := User{
  "Name": "King",
  "Role": "King"
}
if err := userWrapper.New(user); err != nil {
  panic(err)
}
```

Find Many Populate (Using aggregate lookup)

```go
query := bson.D{{
  Key:   "_id",
  Value: mongoose.ObjId(campaignId),
}}

populate := mongoose.Populate{
  As:           "Quizzes.Data",
  ForeignModel: models.QuizTemplates,
  LocalField:   "Quizzes.Ids",
}

campaigns := []Campaign{}
if err := campaignWrapper.FindManyPopulate(query, populate, &campaigns); err != nil {
  utils.Panic(401, "[1] Campaign Not Found", err)
}
```

#### Other utility functions

Easily create `primitives.ObjectId`s from the same object id's stored as string

```go
primitive := mongoose.ObjId("62d04ca380bfc30a1436f590")
```

Get `primitive.DateTime` for storing all your `createdAt` to `time.Now()`

```go
user := User{
  CreatedAt: mongoose.Now() // This will properly create a primitive.DateTime object for the current date time
}

```
