package mongo

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/claudioluciano/goutils/errors"
	"github.com/claudioluciano/goutils/logger"
)

type Client struct {
	mongoClient    *mongo.Client
	logger         *logger.Client
	errors         *errors.Client
	databaseName   string
	collectionName string
}

type NewClientOptions struct {
	Host           string
	Port           string
	User           string
	Password       string
	DatabaseName   string
	CollectionName string
	Logger         *logger.Client
	Errors         *errors.Client
}

func NewClient(opts *NewClientOptions) (*Client, error) {
	ctx := context.TODO()
	clientOpts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", opts.User, opts.Password, opts.Host, opts.Port))

	mongoClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	collectionName := toSnakeCase(opts.CollectionName)

	opts.Logger.Info("connected to mongodb", "db_name", opts.DatabaseName, "db_collection", collectionName)

	return &Client{
		mongoClient:    mongoClient,
		logger:         opts.Logger,
		databaseName:   opts.DatabaseName,
		collectionName: collectionName,
	}, nil
}

func toSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func (c *Client) database() *mongo.Database {
	return c.mongoClient.Database(c.databaseName)
}

func (c *Client) collection() *mongo.Collection {
	db := c.database()
	return db.Collection(c.collectionName)
}
