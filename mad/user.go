package mad

import (
	"context"
	"fmt"
	"github.com/hyaku-roku-ju/til/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Db *mongo.Database
}

func (self *UserRepository) SetPreferredTime(ctx context.Context, id string, preferredTime user.PreferredTime) error {
	collection := self.Db.Collection("users")
	fail := make(chan error)

	go func() {
		id, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			fail <- err
		}

		_, err = collection.UpdateOne(
			ctx,
			bson.M{"_id": id},
			bson.D{{Key: "$set", Value: bson.M{"preferredTime": bson.M{
				"hour": preferredTime.Hour,
				"min":  preferredTime.Min,
			}}}},
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

func (self *UserRepository) Create(ctx context.Context, preferredTime user.PreferredTime) (string, error) {
	collection := self.Db.Collection("users")
	fail := make(chan error)
	success := make(chan string)

	go func() {
		result, err := collection.InsertOne(ctx, bson.M{"preferredTime": bson.M{
			"hour": preferredTime.Hour,
			"min":  preferredTime.Min,
		}})

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
