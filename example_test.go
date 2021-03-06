package lungo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func Example() {
	type post struct {
		Title string `bson:"title"`
	}

	// prepare options
	opts := Options{
		Store: NewMemoryStore(),
	}

	// open database
	client, engine, err := Open(nil, opts)
	if err != nil {
		panic(err)
	}

	// ensure engine is closed
	defer engine.Close()

	// get db
	foo := client.Database("foo")

	// get collection
	bar := foo.Collection("bar")

	// insert post
	_, err = bar.InsertOne(nil, &post{
		Title: "Hello World!",
	})
	if err != nil {
		panic(err)
	}

	// query posts
	csr, err := bar.Find(nil, bson.M{})
	if err != nil {
		panic(err)
	}

	// decode posts
	var posts []post
	err = csr.All(nil, &posts)
	if err != nil {
		panic(err)
	}

	// print documents
	fmt.Printf("%+v", posts)

	// Output:
	// [{Title:Hello World!}]
}
