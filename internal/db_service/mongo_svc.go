package db_service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrNotFound = errors.New("document not found")
	ErrConflict = errors.New("document conflict")
)

type DbService[DocType any] interface {
	CreateDocument(ctx context.Context, id string, document *DocType) error
	FindDocument(ctx context.Context, id string) (*DocType, error)
	FindAllDocuments(ctx context.Context) ([]DocType, error)
	UpdateDocument(ctx context.Context, id string, document *DocType) error
	DeleteDocument(ctx context.Context, id string) error
	Disconnect(ctx context.Context) error
}

type MongoServiceConfig struct {
	Server     string
	DbName     string
	Collection string
	Timeout    time.Duration
}

type mongoSvc[DocType any] struct {
	MongoServiceConfig
	client *mongo.Client
}

func NewMongoService[DocType any](config MongoServiceConfig) DbService[DocType] {
	if config.Server == "" {
		config.Server = "localhost:27017"
	}
	if config.DbName == "" {
		config.DbName = "davgus-ambulance-wl" // specific to this project
	}
	if config.Collection == "" {
		config.Collection = "department" // specific to this project
	}
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}

	return &mongoSvc[DocType]{
		MongoServiceConfig: config,
	}
}

func (m *mongoSvc[DocType]) connect(ctx context.Context) (*mongo.Client, error) {
	if m.client != nil {
		return m.client, nil
	}

	dbHost := os.Getenv("AMBULANCE_API_MONGODB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("AMBULANCE_API_MONGODB_PORT")
	if dbPort == "" {
		dbPort = "27017"
	}
	dbUser := os.Getenv("AMBULANCE_API_MONGODB_USERNAME")
	dbPass := os.Getenv("AMBULANCE_API_MONGODB_PASSWORD")

	uri := fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	if dbUser != "" && dbPass != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPass, dbHost, dbPort)
	}

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}
	m.client = client
	return client, nil
}

func (m *mongoSvc[DocType]) Disconnect(ctx context.Context) error {
	if m.client != nil {
		return m.client.Disconnect(ctx)
	}
	return nil
}

func (m *mongoSvc[DocType]) CreateDocument(ctx context.Context, id string, document *DocType) error {
	ctx, contextCancel := context.WithTimeout(ctx, m.Timeout)
	defer contextCancel()
	client, err := m.connect(ctx)
	if err != nil {
		return err
	}
	db := client.Database(m.DbName)
	collection := db.Collection(m.Collection)
	result := collection.FindOne(ctx, bson.D{{Key: "id", Value: id}})
	switch result.Err() {
	case nil: // no error means there is conflicting document
		return ErrConflict
	case mongo.ErrNoDocuments:
		// do nothing, this is expected
	default: // other errors - return them
		return result.Err()
	}

	_, err = collection.InsertOne(ctx, document)
	return err
}

func (m *mongoSvc[DocType]) FindDocument(ctx context.Context, id string) (*DocType, error) {
	ctx, contextCancel := context.WithTimeout(ctx, m.Timeout)
	defer contextCancel()
	client, err := m.connect(ctx)
	if err != nil {
		return nil, err
	}
	db := client.Database(m.DbName)
	collection := db.Collection(m.Collection)
	result := collection.FindOne(ctx, bson.D{{Key: "id", Value: id}})
	switch result.Err() {
	case nil:
	case mongo.ErrNoDocuments:
		return nil, ErrNotFound
	default: // other errors - return them
		return nil, result.Err()
	}
	var document *DocType
	if err := result.Decode(&document); err != nil {
		return nil, err
	}
	return document, nil
}

func (m *mongoSvc[DocType]) FindAllDocuments(ctx context.Context) ([]DocType, error) {
	ctx, contextCancel := context.WithTimeout(ctx, m.Timeout)
	defer contextCancel()
	client, err := m.connect(ctx)
	if err != nil {
		return nil, err
	}
	db := client.Database(m.DbName)
	collection := db.Collection(m.Collection)
	
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var documents []DocType
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, err
	}
	if documents == nil {
		documents = []DocType{}
	}
	return documents, nil
}

func (m *mongoSvc[DocType]) UpdateDocument(ctx context.Context, id string, document *DocType) error {
	ctx, contextCancel := context.WithTimeout(ctx, m.Timeout)
	defer contextCancel()
	client, err := m.connect(ctx)
	if err != nil {
		return err
	}
	db := client.Database(m.DbName)
	collection := db.Collection(m.Collection)
	result := collection.FindOne(ctx, bson.D{{Key: "id", Value: id}})
	switch result.Err() {
	case nil:
	case mongo.ErrNoDocuments:
		return ErrNotFound
	default: // other errors - return them
		return result.Err()
	}
	_, err = collection.ReplaceOne(ctx, bson.D{{Key: "id", Value: id}}, document)
	return err
}

func (m *mongoSvc[DocType]) DeleteDocument(ctx context.Context, id string) error {
	ctx, contextCancel := context.WithTimeout(ctx, m.Timeout)
	defer contextCancel()
	client, err := m.connect(ctx)
	if err != nil {
		return err
	}
	db := client.Database(m.DbName)
	collection := db.Collection(m.Collection)
	result := collection.FindOne(ctx, bson.D{{Key: "id", Value: id}})
	switch result.Err() {
	case nil:
	case mongo.ErrNoDocuments:
		return ErrNotFound
	default: // other errors - return them
		return result.Err()
	}
	_, err = collection.DeleteOne(ctx, bson.D{{Key: "id", Value: id}})
	return err
}
