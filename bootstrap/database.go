package bootstrap

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/mongo"
)

func NewMongoDatabase(env *Env) mongo.Client {
	dbHost := env.DBHost
	dbPort := env.DBPort
	dbUser := env.DBUser
	dbPass := env.DBPass

	mongodbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPass, dbHost, dbPort)

	if dbUser == "" || dbPass == "" {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	}

	const (
		maxAttempts = 30
		retryDelay  = 2 * time.Second
		pingTimeout = 2 * time.Second
	)

	var lastErr error
	for i := 0; i < maxAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
		client, err := mongo.NewClient(mongodbURI)
		if err == nil {
			err = client.Ping(ctx)
		}
		cancel()

		if err == nil {
			return client
		}

		if client != nil {
			_ = client.Disconnect(context.TODO())
		}

		lastErr = err
		log.Printf("MongoDB not ready (attempt %d/%d): %v", i+1, maxAttempts, err)
		time.Sleep(retryDelay)
	}

	log.Fatal(lastErr)
	return nil
}

func CloseMongoDBConnection(client mongo.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}
