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

type TelegramIndex struct {
	name      string
	isUnique  bool
	direction int
}

var telegramIndexes = []TelegramIndex{
	{name: "telegramId", isUnique: true, direction: 1},
	{name: "userId", isUnique: true, direction: 1},
}

// Get the Telegram collection and create indexes if they have not yet been created
func NewTelegramRepository(ctx context.Context, db *mongo.Database) TelegramRepository {
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

	// create a map for easy look up
	indexMap := make(map[string]bool)
	for _, indexBson := range indexes {
		var indexName string
		for key, value := range indexBson {
			name, isString := value.(string)
			if key == "name" && isString {
				indexName = name
			}
		}

		indexMap[indexName] = true
	}

	// map through the telegramIndexes to check if the index has been created or not
	for _, tIndex := range telegramIndexes {
		// assume the index hasn't been created if not in the indexMap
		if _, ok := indexMap[tIndex.name]; !ok {
			// create mongo index options
			opts := mongo.IndexModel{
				Keys:    bson.M{tIndex.name: tIndex.direction},
				Options: options.Index().SetUnique(tIndex.isUnique).SetName(tIndex.name),
			}
			_, err := col.Indexes().CreateOne(ctx, opts)

			if err != nil {
				log.Fatal("Unable to create telegram index", err)
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
