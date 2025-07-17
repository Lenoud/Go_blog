package mongodb

import (
	"context"
	"log"
	"time"

	"blog/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.GlobalConfig.MongoDB.URI)
	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// 测试连接
	err = Client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	Database = Client.Database(config.GlobalConfig.MongoDB.Database)
	log.Println("MongoDB connected successfully")
	return nil
}
