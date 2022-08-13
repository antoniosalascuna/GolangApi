package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectionDB() *mongo.Client {

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(("mongodb://mongo:27017/PR_collection")))

	if err != nil {
		log.Fatal(err)
	}
	///Checkea que la conexion se esta estableciendo con el servidor de mongo(container)
	err = client.Ping(ctx, nil)

	if err != nil {
		fmt.Println("Error To connect MongoDb")
		log.Fatal(err)
	}

	fmt.Println("Connect To MongoDb")

	return client
}

func Disconnect(client *mongo.Client, context context.Context) {

	defer client.Disconnect(context)

	fmt.Println("Disconnect To MongoDb")

}
