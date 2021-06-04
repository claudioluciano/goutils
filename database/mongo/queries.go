package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (c *Client) FindOneByID(ctx context.Context, id string, target interface{}) error {
	if err := c.collection().FindOne(ctx, &bson.M{"_id": id}).Decode(target); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.errors.NotFound("error", err)
		}

		return c.errors.Internal(err)
	}

	return nil
}

func (c *Client) FindOne(ctx context.Context, filter interface{}, target interface{}, opts ...*options.FindOneOptions) error {
	if err := c.collection().FindOne(ctx, &filter, opts...).Decode(target); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.errors.NotFound("error", err)
		}

		return c.errors.Internal(err)
	}

	return nil
}

func (m *Client) FindPaginated(ctx context.Context, filter interface{}, skip int64, limit int64) (*mongo.Cursor, error) {
	opts := options.Find().SetSkip(skip).SetLimit(limit)
	cur, err := m.collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, m.errors.Internal(err)
	}

	return cur, err
}

func (m *Client) FindMany(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	cur, err := m.collection().Find(ctx, filter, opts...)
	if err != nil {
		return nil, m.errors.Internal(err)
	}

	return cur, err
}

func (c *Client) InsertOne(ctx context.Context, data interface{}) error {
	_, err := c.collection().InsertOne(ctx, &data)
	if err != nil {
		return c.errors.Internal(err)
	}

	return nil
}

func (c *Client) InsertMany(ctx context.Context, data []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	result, err := c.collection().InsertMany(ctx, data, opts...)
	if err != nil {
		return nil, c.errors.Internal(err)
	}

	return result, nil
}

func (c *Client) UpdateOneByID(ctx context.Context, id string, update interface{}, target interface{}, opts ...*options.UpdateOptions) error {
	dbOpts := options.FindOneAndUpdate().SetReturnDocument(1)
	if err := c.findOneAndUpdate(ctx, &bson.M{"_id": id}, update, dbOpts).Decode(target); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.errors.NotFound("error", err)
		}

		return c.errors.Internal(err)
	}

	return nil
}

func (c *Client) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	result, err := c.collection().UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return nil, c.errors.Internal(err)
	}

	return result, nil
}

func (m *Client) findOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	result := m.collection().FindOneAndUpdate(ctx, filter, update, opts...)

	return result
}

func (m *Client) ReplaceOneByID(ctx context.Context, id string, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	result, err := m.collection().ReplaceOne(ctx, &bson.M{"_id": id}, replacement, opts...)
	if err != nil {
		return nil, m.errors.Internal(err)
	}

	return result, nil
}

func (m *Client) DeleteOneByID(ctx context.Context, id string, target interface{}, opts ...*options.FindOneAndDeleteOptions) error {
	if err := m.collection().FindOneAndDelete(ctx, &bson.M{"_id": id}, opts...).Decode(target); err != nil {
		if err == mongo.ErrNoDocuments {
			return m.errors.NotFound("error", err)
		}

		return m.errors.Internal(err)
	}

	return nil
}

func (m *Client) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	result, err := m.collection().DeleteMany(ctx, filter, opts...)
	if err != nil {
		return nil, m.errors.Internal(err)
	}

	return result, nil
}
