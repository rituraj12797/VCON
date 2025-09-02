package db

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func DBConnect() {

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error : ", err)
		panic(err)
	}

	mongouri := os.Getenv("MONGO_URI")

	fmt.Println(" this is our mongo connection url : ", mongouri)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(mongouri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)

	database := client.Database("test")
	// collection := database.Collection("document")
	err = database.CreateCollection(context.TODO(),"content_sentence")

	if err != nil {
		panic(err)
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

}
