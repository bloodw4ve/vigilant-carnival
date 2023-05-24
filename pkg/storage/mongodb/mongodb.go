package mongodb

import (
	"GoNews/pkg/storage"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	Client *mongo.Client
}

const (
	dbName         = "gonews"
	collectionName = "posts"
)

func New(constr string) (*Store, error) {
	mongoOpts := options.Client().ApplyURI(constr)
	client, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	s := Store{
		Client: client,
	}
	return &s, nil
}

func (s *Store) Close() {
	s.Client.Disconnect(context.Background())
}

func (s *Store) Posts() ([]storage.Post, error) {
	collection := s.Client.Database(dbName).Collection(collectionName)
	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)

	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())
	var data []storage.Post
	for cur.Next(context.Background()) {
		var p storage.Post
		err := cur.Decode(&p)

		if err != nil {
			return nil, err
		}

		data = append(data, p)
	}

	return data, cur.Err()
}

func (s *Store) AddPost(post storage.Post) error {
	collection := s.Client.Database(dbName).Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), post)

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdatePost(post storage.Post) error {
	collection := s.Client.Database(dbName).Collection(collectionName)
	filter := bson.M{"id": post.ID}
	update := bson.M{"$set": post}
	_, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) DeletePost(post storage.Post) error {
	collection := s.Client.Database(dbName).Collection(collectionName)
	filter := bson.M{"id": post.ID}
	_, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}
	return nil
}
