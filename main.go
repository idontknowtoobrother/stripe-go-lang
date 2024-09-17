package main

import (
	"context"
	"fmt"

	"github.com/gin-contrib/cors"
	controllers "github.com/idontknowtoobrother/stripe-go-lang/Controllers"
	repository "github.com/idontknowtoobrother/stripe-go-lang/Repository"
	routes "github.com/idontknowtoobrother/stripe-go-lang/Routes"
	utils "github.com/idontknowtoobrother/stripe-go-lang/Utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(utils.GetEnv("MONGODB_URI")))
	if err != nil {
		panic(err)
	}
	defer mongoClient.Disconnect(ctx)

	repo := repository.NewRepo(ctx, mongoClient.Database(utils.GetEnv("DB_NAME")))
	controller := controllers.NewController(ctx, repo)

	engine := routes.SetupRoutes(controller)
	engine.Use(cors.Default())

	fmt.Println("Server is running on port: ", utils.GetEnv("PORT"))
	if err := engine.Run(":" + utils.GetEnv("PORT")); err != nil {
		panic(err)
	}
}
