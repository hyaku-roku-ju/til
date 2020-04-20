package mad

import (
	"context"
	"fmt"

	"github.com/hyaku-roku-ju/til/learning"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LearningRepository struct {
	Db *mongo.Database
}

var indexes = []mongo.IndexModel{
	{
		Keys: bson.D{
			{Key: "reporterId", Value: 1},
			{Key: "confirmed", Value: 1},
		},
		Options: options.Index().SetBackground(true).SetName("reporterId.confirmed"),
	},
}

func NewLearningRepository(ctx context.Context, db *mongo.Database) (LearningRepository, error) {
	collection := db.Collection("learnings")
	fail := make(chan error)

	go func() {
		var existingIndexes []bson.M

		indexView := collection.Indexes()
		cursor, err := indexView.List(ctx)
		if err != nil {
			fail <- err
		}
		err = cursor.All(ctx, &existingIndexes)
		if err != nil {
			fail <- err
		}

		existingIndexSet := make(map[string]bool)
		for _, indexDocument := range existingIndexes {
			indexName, ok := indexDocument["name"].(string)
			if !ok {
				fail <- fmt.Errorf("Failed to get name field from existing mongo index while creating Learning repository")
			}
			existingIndexSet[indexName] = true
		}

		indexesToCreate := make([]mongo.IndexModel, 0)
		for _, indexDocument := range indexes {
			if _, exists := existingIndexSet[*indexDocument.Options.Name]; exists {
				continue
			}
			// index does not exist
			indexesToCreate = append(indexesToCreate, indexDocument)
		}
		// nothing to create, exit early
		if len(indexesToCreate) == 0 {
			fail <- nil
		}
		_, err = indexView.CreateMany(ctx, indexesToCreate)
		fail <- err
	}()

	select {
	case <-ctx.Done():
		return LearningRepository{}, ctx.Err()
	case err := <-fail:
		if err != nil {
			return LearningRepository{}, err
		} else {
			return LearningRepository{Db: db}, nil
		}
	}
}

func (self *LearningRepository) GetConfirmedLearning(ctx context.Context, reporterId string, skip int) (learning.Learning, error) {
	collection := self.Db.Collection("learnings")
	fail := make(chan error)
	success := make(chan learning.Learning)

	go func() {
		reporterId, err := primitive.ObjectIDFromHex(reporterId)
		if err != nil {
			fail <- err
		}
		opts := options.FindOne()
		// first learning would be {skip: 0},
		opts.SetSkip(int64(skip))
		var randomLearning learning.Learning
		err = collection.FindOne(
			ctx,
			bson.D{
				{Key: "reporterId", Value: reporterId},
				{Key: "confirmed", Value: true},
			},
			opts,
		).Decode(&randomLearning)

		if err != nil {
			fail <- err
		}

		success <- randomLearning
	}()

	select {
	case <-ctx.Done():
		return learning.Learning{}, ctx.Err()
	case err := <-fail:
		return learning.Learning{}, err
	case randomLearning := <-success:
		return randomLearning, nil
	}
}

func (self *LearningRepository) CountConfirmedLearnings(ctx context.Context, reporterId string) (int, error) {
	collection := self.Db.Collection("learnings")
	fail := make(chan error)
	success := make(chan int64)

	go func() {
		reporterId, err := primitive.ObjectIDFromHex(reporterId)
		if err != nil {
			fail <- err
		}

		count, err := collection.CountDocuments(ctx, bson.D{
			{Key: "reporterId", Value: reporterId},
			{Key: "confirmed", Value: true},
		})

		if err != nil {
			fail <- err
		}

		success <- count
	}()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case err := <-fail:
		return 0, err
	case count := <-success:
		// if somebody has more than 2^32 learnings
		// it would overflow on 32bit systems
		return int(count), nil
	}
}

func (self *LearningRepository) ConfirmLearning(ctx context.Context, id string) error {
	collection := self.Db.Collection("learnings")
	fail := make(chan error)

	go func() {
		id, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			fail <- err
		}

		_, err = collection.UpdateOne(
			ctx,
			bson.M{"_id": id},
			bson.D{
				{Key: "$set", Value: bson.M{"confirmed": true}},
			},
		)

		fail <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-fail:
		return err
	}
}

func (self *LearningRepository) StoreLearning(ctx context.Context, learning learning.Learning) (string, error) {
	collection := self.Db.Collection("learnings")
	fail := make(chan error)
	success := make(chan string)

	go func() {
		id, err := primitive.ObjectIDFromHex(learning.Id)
		if err != nil {
			fail <- err
		}
		reporterId, err := primitive.ObjectIDFromHex(learning.ReporterId)
		if err != nil {
			fail <- err
		}
		learningToInsert := bson.M{
			"_id":         id,
			"description": learning.Description,
			"topics":      learning.Topics,
			"reporterId":  reporterId,
			"confirmed":   learning.Confirmed,
		}
		result, err := collection.InsertOne(ctx, learningToInsert)

		if err != nil {
			fail <- err
		}
		objectId, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			fail <- fmt.Errorf("Failed to cast insertOne result to ObjectID")
		}
		success <- objectId.Hex()
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
