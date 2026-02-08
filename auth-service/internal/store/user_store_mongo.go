package store

import (
	"context"
	"time"

	"auth-service/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserStore struct {
	users    *mongo.Collection
	counters *mongo.Collection
}

func NewMongoUserStore(db *mongo.Database) *MongoUserStore {
	s := &MongoUserStore{
		users:    db.Collection("users"),
		counters: db.Collection("counters"),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = s.EnsureIndexes(ctx)
	return s
}

func (s *MongoUserStore) EnsureIndexes(ctx context.Context) error {
	_, err := s.users.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	return err
}

func (s *MongoUserStore) NextUserID(ctx context.Context) (int64, error) {
	res := s.counters.FindOneAndUpdate(
		ctx,
		bson.M{"_id": "users"},
		bson.M{"$inc": bson.M{"value": 1}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	)

	var c model.CounterDoc
	if err := res.Decode(&c); err != nil {
		return 0, err
	}
	return c.Value, nil
}

func (s *MongoUserStore) InsertUser(ctx context.Context, u model.UserDoc) error {
	_, err := s.users.InsertOne(ctx, u)
	return err
}

func (s *MongoUserStore) FindByUsername(ctx context.Context, username string) (model.UserDoc, bool, error) {
	var u model.UserDoc
	err := s.users.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return model.UserDoc{}, false, nil
	}
	if err != nil {
		return model.UserDoc{}, false, err
	}
	return u, true, nil
}
