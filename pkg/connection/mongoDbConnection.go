package connection

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"sync"
)

//Initialized and exposed through  GetMongoClient().
var clientInstance *mongo.Client
//Used during creation of singleton client object in GetMongoClient().
var clientInstanceError error
//Used to execute client creation procedure only once.
var mongoOnce sync.Once
//I have used below constants just to hold required database config's.
var (
	CONNECTIONSTRING = os.Getenv("Connection_String")
	DB               = os.Getenv("DB")
	ORDERS           = os.Getenv("Col")
)
//GetMongoClient - Return mongodb connection to work with
func GetMongoClient() (*mongo.Client, error) {
	//Perform connection creation operation only once.
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)
		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}
		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = client
	})
	return clientInstance, clientInstanceError
}
