package db

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBConnect() (*mongo.Database, error) {

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error : ", err)
		panic(err)
	}

	mongouri := os.Getenv("MONGO_URI")

	// fmt.Println(" this is our mongo connection url : ", mongouri)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(mongouri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)

	database := client.Database("test")
	// collection := database.Collection("document")

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return database, nil
}
