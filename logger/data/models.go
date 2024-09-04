package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var client *mongo.Client

// Models is the type for this package
type Models struct {
	LogEntry LogEntry
}

// LogEntry stores the log entry in the database
type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreateAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// New will create a new log entry model
func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

// Insert will insert a new log entry
func (l *LogEntry) Insert(entry LogEntry) error {
	// create a collection
	collection := client.Database("logs").Collection("logs")

	// insert the entry
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreateAt:  time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into logs: ", err)
		return err
	}

	return nil
}

// All will return all log entries
func (l *LogEntry) All() ([]*LogEntry, error) {
	// create a context
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// create a collection
	collection := client.Database("logs").Collection("logs")

	// create options
	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	// find all logs
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Error finding all logs: ", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	// create an array of logs
	var logs []*LogEntry

	// loop through the logs
	for cursor.Next(ctx) {
		var item LogEntry
		err = cursor.Decode(&item)
		if err != nil {
			log.Println("Error decoding log: ", err)
			return nil, err
		}

		logs = append(logs, &item)
	}

	return logs, nil
}

// GetOne will return one log entry
func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	// create a context
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// create a collection
	collection := client.Database("logs").Collection("logs")

	// convert id to object id
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// find the log
	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docId}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// DropCollection will drop the logs collection
func (l *LogEntry) DropCollection() error {
	// create a context
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// create a collection
	collection := client.Database("logs").Collection("logs")

	// drop the collection
	err := collection.Drop(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Update will update one log
func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	// create a context
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// create a collection
	collection := client.Database("logs").Collection("logs")

	// get doc id
	docId, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	// update the log
	result, err := collection.UpdateOne(ctx, bson.M{"_id": docId}, bson.D{
		{"$set", bson.D{
			{"name", l.Name},
			{"data", l.Data},
			{"updated_at", time.Now()},
		}},
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
