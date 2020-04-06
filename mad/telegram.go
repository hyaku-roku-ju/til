package mad

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type TelegramRepository struct {
	Db *mongo.Database
}

var indexesOpts = []mongo.IndexModel{
	{Keys: bson.M{"telegramId": 1}, Options: options.Index().SetUnique(true)},
	{Keys: bson.M{"userId": 1}, Options: options.Index().SetUnique(true)},
}

// Get the Telegram collection and create indexes if they have not yet been
func CreateCollection(ctx context.Context, db *mongo.Database) TelegramRepository {
	col := db.Collection("telegram")

	opts := options.ListIndexes()
	cursor, err := col.Indexes().List(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	var indexes []bson.M
	if err = cursor.All(ctx, &indexes); err != nil {
		log.Fatal(err)
	}

	// Assume that 1 index would be _id
	if len(indexes) <= 1 {
		// Set all the indexes that are unique
		for _, opts := range indexesOpts {
			_, err := col.Indexes().CreateOne(ctx, opts)

			if err != nil {
				// Is this a fatal error?
				log.Fatal(err)
			}
		}
	}

	return TelegramRepository{Db: db}
}

func (self *TelegramRepository) Create(ctx context.Context, userId string, telegramId string) error {
	collection := self.Db.Collection("telegram")
	fail := make(chan error)

	go func() {
		_, err := collection.InsertOne(ctx, bson.M{"userId": userId, "telegramId": telegramId})

		fail <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-fail:
		return err
	}
}

func (self *TelegramRepository) GetIdByUserId(ctx context.Context, userId string) (string, error) {
	collection := self.Db.Collection("telegram")
	fail := make(chan error)
	success := make(chan string)

	go func() {
		var result bson.M
		err := collection.FindOne(ctx, bson.M{"userId": userId}).Decode(&result)
		if err != nil {
			fail <- err
		}

		if telegramId, ok := result["telegramId"].(string); ok {
			success <- telegramId
		} else {
			fail <- fmt.Errorf("Unable to get telegramId from result")
		}
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-fail:
		return "", err
	case id := <-success:
		return id, nil
	}
}
