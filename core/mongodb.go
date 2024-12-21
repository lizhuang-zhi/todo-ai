package core

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Mongodb struct {
	url      string
	database string
	client   *mongo.Client
}

func NewMongoDB(url string) (*Mongodb, error) {
	cs, err := connstring.Parse(url)
	if err != nil {
		return nil, err
	}

	client, err := mongo.Connect(context.Background(), mongoOptions.Client().ApplyURI(url).SetConnectTimeout(10*time.Second))
	if err != nil {
		return nil, err
	}

	db := &Mongodb{
		url:      url,
		database: cs.Database,
		client:   client,
	}

	return db, nil
}

func (db *Mongodb) C(collection string) *mongo.Collection {
	return db.client.Database(db.database).Collection(collection)
}

func (db *Mongodb) Ping(ctx context.Context) error {
	return db.client.Ping(ctx, readpref.Primary())
}

func (db *Mongodb) InsertOne(ctx context.Context, collection string, doc interface{}) (id interface{}, err error) {
	var result *mongo.InsertOneResult

	result, err = db.C(collection).InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (db *Mongodb) Insert(ctx context.Context, collection string, docs []interface{}) (ids []interface{}, err error) {
	var result *mongo.InsertManyResult

	result, err = db.C(collection).InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}

	return result.InsertedIDs, nil
}

func (db *Mongodb) DeleteOne(ctx context.Context, collection string, filter interface{}) (err error) {
	if filter == nil {
		filter = bson.M{}
	}

	_, err = db.C(collection).DeleteOne(ctx, filter)
	return err
}

func (db *Mongodb) Delete(ctx context.Context, collection string, filter interface{}) (cnt int64, err error) {
	if filter == nil {
		filter = bson.M{}
	}

	var result *mongo.DeleteResult

	result, err = db.C(collection).DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (db *Mongodb) UpdateOne(ctx context.Context, collection string, filter interface{}, update interface{}) (modified bool, err error) {
	if filter == nil {
		filter = bson.M{}
	}

	if _, ok := update.(bson.M); !ok {
		update = bson.M{"$set": update}
	}

	var result *mongo.UpdateResult
	result, err = db.C(collection).UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}

	return result.ModifiedCount > 0, nil
}

func (db *Mongodb) Update(ctx context.Context, collection string, filter interface{}, update interface{}) (cnt int64, err error) {
	if filter == nil {
		filter = bson.M{}
	}

	if _, ok := update.(bson.M); !ok {
		update = bson.M{"$set": update}
	}

	var result *mongo.UpdateResult
	result, err = db.C(collection).UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

func (db *Mongodb) Upsert(ctx context.Context, collection string, filter interface{}, update interface{}) (upsert bool, id interface{}, err error) {
	if filter == nil {
		filter = bson.M{}
	}

	if _, ok := update.(bson.M); !ok {
		update = bson.M{"$set": update}
	}

	var result *mongo.UpdateResult
	result, err = db.C(collection).UpdateOne(ctx, filter, update, mongoOptions.Update().SetUpsert(true))
	if err != nil {
		return false, nil, err
	}

	return result.UpsertedCount > 0, result.UpsertedID, nil
}

func (db *Mongodb) FindOne(ctx context.Context, collection string, filter interface{}, result interface{}) (err error) {
	if filter == nil {
		filter = bson.M{}
	}

	err = db.C(collection).FindOne(ctx, filter).Decode(result)
	return err
}

func (db *Mongodb) Find(ctx context.Context, collection string, filter interface{}, cols interface{}, results interface{}) (err error) {
	if filter == nil {
		filter = bson.M{}
	}

	opts := mongoOptions.Find()
	opts.SetProjection(cols)

	var cursor *mongo.Cursor
	cursor, err = db.C(collection).Find(ctx, filter, opts)
	if err != nil {
		return err
	}

	return cursor.All(ctx, results)
}

func (db *Mongodb) FindRange(ctx context.Context, collection string, filter interface{}, skip, limit int64, results interface{}) (err error) {
	if filter == nil {
		filter = bson.M{}
	}

	opts := mongoOptions.Find()

	opts.SetSkip(skip)
	opts.SetLimit(limit)

	var cursor *mongo.Cursor
	cursor, err = db.C(collection).Find(ctx, filter, opts)
	if err != nil {
		return err
	}

	return cursor.All(ctx, results)
}

func (db *Mongodb) FindSortDesRange(ctx context.Context, collection string, filter interface{}, sortField string, skip, limit int64, results interface{}) (err error) {
	if filter == nil {
		filter = bson.M{}
	}

	opts := mongoOptions.Find()

	opts.SetSort(bson.D{{Key: sortField, Value: -1}})
	opts.SetSkip(skip)
	opts.SetLimit(limit)

	var cursor *mongo.Cursor
	cursor, err = db.C(collection).Find(ctx, filter, opts)
	if err != nil {
		return err
	}

	return cursor.All(ctx, results)
}

func (db *Mongodb) FindSortAscRange(ctx context.Context, collection string, filter interface{}, sortField string, skip, limit int64, results interface{}) (err error) {
	if filter == nil {
		filter = bson.M{}
	}

	opts := mongoOptions.Find()

	opts.SetSort(bson.D{{Key: sortField, Value: 1}})
	opts.SetSkip(skip)
	opts.SetLimit(limit)

	var cursor *mongo.Cursor
	cursor, err = db.C(collection).Find(ctx, filter, opts)
	if err != nil {
		return err
	}

	return cursor.All(ctx, results)
}

func (db *Mongodb) FindOneAndUpdate(ctx context.Context, collection string, filter, update, result interface{}) (err error) {
	if filter == nil {
		filter = bson.M{}
	}

	return db.C(collection).FindOneAndUpdate(ctx, filter, update).Decode(result)
}

func (db *Mongodb) Count(ctx context.Context, collection string, filter interface{}) (cnt int64, err error) {
	if filter == nil {
		filter = bson.M{}
	}

	cnt, err = db.C(collection).CountDocuments(ctx, filter)
	return cnt, err
}

// 批量操作
func (db *Mongodb) BulkWrite(ctx context.Context, collection string, models []mongo.WriteModel) (ids map[int64]interface{}, err error) {
	var result *mongo.BulkWriteResult

	result, err = db.C(collection).BulkWrite(ctx, models)
	if err != nil {
		return nil, err
	}

	return result.UpsertedIDs, nil
}

// 索引
type Index struct {
	Unique bool
	Fields []string
}

// 创建索引
func NewIndex(unique bool, fields ...string) *Index {
	return &Index{
		Unique: unique,
		Fields: fields,
	}
}

// 创建多个索引
func (db *Mongodb) CreateIndexes(ctx context.Context, collection string, indexes ...*Index) error {
	for _, index := range indexes {
		if err := db.CreateIndex(ctx, collection, index.Fields, index.Unique); err != nil {
			return err
		}
	}
	return nil
}

// 创建索引
func (db *Mongodb) CreateIndex(ctx context.Context, collection string, fields []string, unique bool) error {
	keys := bson.D{}
	for _, field := range fields {
		keys = append(keys, bson.E{Key: field, Value: 1})
	}

	_, err := db.C(collection).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    keys,
		Options: mongoOptions.Index().SetUnique(unique),
	})

	if err == nil {
		return nil
	}

	return fmt.Errorf("%v: %w", collection, err)
}

func (db *Mongodb) ListIndex(ctx context.Context, collection string) ([]bson.M, error) {
	cursor, err := db.C(collection).Indexes().List(ctx)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	err = cursor.All(context.TODO(), &results)
	return results, err
}

func (db *Mongodb) ListCollectionNames(ctx context.Context, filter interface{}) ([]string, error) {
	return db.client.Database(db.database).ListCollectionNames(ctx, filter)
}

func (db *Mongodb) DropCollection(ctx context.Context, collection string) error {
	return db.C(collection).Drop(ctx)
}

func (db *Mongodb) DropDatabase(ctx context.Context) error {
	return db.client.Database(db.database).Drop(ctx)
}

func (db *Mongodb) Close(ctx context.Context) error {
	return db.client.Disconnect(ctx)
}

func (db *Mongodb) RunCommand(ctx context.Context, runCommand interface{}) (bson.Raw, error) {
	return db.client.Database(db.database).RunCommand(ctx, runCommand).DecodeBytes()
}

func (db *Mongodb) StartSession() (mongo.Session, error) {
	return db.client.StartSession()
}
