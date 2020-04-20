package mad

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	var models []mongo.IndexModel
	for _, tIndex := range telegramIndexes {
		// assume the index hasn't been created if not in the indexMap
		if _, ok := indexMap[tIndex.name]; !ok {
			// create mongo index options
			model := mongo.IndexModel{
				Keys:    bson.M{tIndex.name: tIndex.direction},
				Options: options.Index().SetUnique(tIndex.isUnique).SetName(tIndex.name),
			}
			models = append(models, model)
		}
	}

	if len(models) > 0 {
		opts := options.CreateIndexes()
		_, err := col.Indexes().CreateMany(ctx, models, opts)

		if err != nil {
			log.Fatal("Unable to create telegram indexes", err, len(models))
		}
	}

	return TelegramRepository{Db: db}
}

func (self *TelegramRepository) Create(ctx context.Context, userId string, telegramId string) error {
	collection := self.Db.Collection("telegram")
	fail := make(chan error)

	go func() {
		id, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			fail <- err
		}

		_, err = collection.InsertOne(ctx, bson.M{"userId": id, "telegramId": telegramId})

		fail <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-fail:
		return err
	}
}

func (self *TelegramRepository) GetTelegramIdByUserId(ctx context.Context, userId string) (string, error) {
	collection := self.Db.Collection("telegram")
	fail := make(chan error)
	success := make(chan string)

	go func() {
		var result bson.M

		id, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			fail <- err
		}

		type fields struct {
			TelegramId int `bson:"telegramId"`
		}
		projection := fields{TelegramId: 1}
		opts := options.FindOne().SetProjection(projection)
		err = collection.FindOne(ctx, bson.M{"userId": id}, opts).Decode(&result)
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

	func (self *TelegramRepository) GetUserIdByTelegramId(ctx context.Context, telegramId string) (string, error) {
		collection := self.Db.Collection("telegram")
		fail := make(chan error)
		success := make(chan string)

		go func() {
			var result bson.M

			type fields struct {
				UserId int `bson:"userId"`
			}
			projection := fields{UserId: 1}
			opts := options.FindOne().SetProjection(projection)
			err := collection.FindOne(ctx, bson.M{"telegramId": telegramId}, opts).Decode(&result)
			if err != nil {
				fail <- err
			}

			if objectId, ok := result["userId"].(primitive.ObjectID); ok {
				success <- objectId.Hex()
			} else {
				fail <- fmt.Errorf("Failed to cast result to objectId")
			}
		}()

		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case err := <- fail:
			return "", err
		case id := <- success:
			return id, nil
		}
	}