package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB_PASSWORD = "DB_PASSWORD"

type MongoDB struct {
	Client *mongo.Client
	config *config.Config
}

func NewMongoDB(client *mongo.Client) (types.ISQL, error) {
	return &MongoDB{
		Client: client,
		config: &config.Config{},
	}, nil
}

func NewMongoDBWithConfig(dbConfig *config.Config) (types.ISQL, error) {
	if os.Getenv(DB_PASSWORD) == "" || len(os.Getenv(DB_PASSWORD)) == 0 { // added mysql to be more verbose about the db type
		return nil, fmt.Errorf("please set %s env variable for the database", DB_PASSWORD)
	}

	mongoEndPoint := dbURLMongoDB(dbConfig)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndPoint))
	if err != nil {
		log.Fatal(err)
	}

	return &MongoDB{
		Client: client,
		config: dbConfig,
	}, nil
}

// Schema returns the schema of a collection in mongoDB.
func (m *MongoDB) Schema(Collection string) (types.Table, error) {

	var schema bson.M
	ctx := context.Background()
	err := m.Client.Database(m.config.DatabaseName).Collection(Collection).FindOne(ctx, bson.M{}).Decode(&schema)
	if err != nil {
		return types.Table{}, err
	}

	return types.Table{
		BSON: schema,
	}, nil

}

// Tables returns the tables(collections) in a database in mongoDB.
func (m *MongoDB) Tables(databaseName string) ([]string, error) {
	var collections []string
	ctx := context.Background()
	collections, err := m.Client.Database(databaseName).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return []string{}, err
	}
	return collections, nil

}

// Execute executes a query on mongoDB.
func (m *MongoDB) Execute(query string) ([]byte, error) {
	return []byte{}, nil
}

func dbURLMongoDB(dbConfig *config.Config) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s",
		dbConfig.Username,
		DB_PASSWORD,
		dbConfig.Host,
		dbConfig.Port,
	)
}
