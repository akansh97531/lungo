package lungo

import (
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const testDB = "test-lungo"

var sharedNativeClient IClient
var sharedClient IClient

var collCounter = 0

func init() {
	nativeClient, err := Connect(nil, options.Client().ApplyURI("mongodb://localhost"))
	if err != nil {
		panic(err)
	}

	err = nativeClient.Database("test-lungo").Drop(nil)
	if err != nil {
		panic(err)
	}

	sharedNativeClient = nativeClient

	client2, err := Open(nil, ClientOptions{
		Store: NewMemoryStore(),
	})
	if err != nil {
		panic(err)
	}

	sharedClient = client2
}

func clientTest(t *testing.T, fn func(t *testing.T, c IClient)) {
	t.Run("NativeClient", func(t *testing.T) {
		fn(t, sharedNativeClient)
	})

	t.Run("Client", func(t *testing.T) {
		fn(t, sharedClient)
	})
}

func databaseTest(t *testing.T, fn func(t *testing.T, d IDatabase)) {
	clientTest(t, func(t *testing.T, client IClient) {
		fn(t, client.Database(testDB))
	})
}

func collectionTest(t *testing.T, fn func(t *testing.T, c ICollection)) {
	collCounter++
	name := fmt.Sprintf("n-%d", collCounter)

	clientTest(t, func(t *testing.T, client IClient) {
		fn(t, client.Database(testDB).Collection(name))
	})
}

func readAll(csr ICursor) []bson.M {
	out := make([]bson.M, 0)
	err := csr.All(nil, &out)
	if err != nil {
		panic(err)
	}

	return out
}

func dumpCollection(c ICollection, clean bool) []bson.M {
	csr, err := c.Find(nil, bson.M{})
	if err != nil {
		panic(err)
	}

	out := readAll(csr)

	if clean {
		for _, item := range out {
			delete(item, "_id")
		}
	}

	return out
}