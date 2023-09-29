package migration

import (
	"context"
	"log"

	"github.com/Rayato159/hello-sekai-shop-tutorial/config"
	"github.com/Rayato159/hello-sekai-shop-tutorial/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func paymentDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(pctx, cfg).Database("payment_db")
}

func PaymentMigrate(pctx context.Context, cfg *config.Config) {
	db := paymentDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	col := db.Collection("payment_queue")

	results, err := col.InsertOne(pctx, bson.M{"offset": -1}, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Migrate payment completed: ", results)
}
