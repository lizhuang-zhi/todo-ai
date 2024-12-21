package core

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

type IntUUID interface {
	Get() (int64, error)
	Init(int64) error
}

const (
	onceNum = 1
)

type count struct {
	Name  string `bson:"name"`
	Count int64  `bson:"count"`
}

type mongoProvider struct {
	sync.Mutex
	mongo      *Mongodb
	collection string
	name       string
	currID     int64
	maxID      int64
}

func NewUUID(mongo *Mongodb, collection, name string) IntUUID {
	return &mongoProvider{
		mongo:      mongo,
		collection: collection,
		name:       name,
	}
}

func (a *mongoProvider) Init(beginID int64) error {
	if err := a.mongo.CreateIndex(context.Background(), a.collection, []string{"name"}, true); err != nil {
		return err
	}
	_, _ = a.mongo.InsertOne(context.Background(), a.collection, count{Name: a.name, Count: beginID})
	return nil
}

func (a *mongoProvider) Get() (int64, error) {
	a.Lock()
	defer a.Unlock()

	if a.currID == a.maxID {
		res := &count{}
		err := a.mongo.FindOneAndUpdate(context.Background(), a.collection, bson.M{"name": a.name}, bson.M{"$inc": bson.M{"count": onceNum}}, res)
		if err != nil {
			return 0, err
		}
		a.currID = res.Count + 1
		a.maxID = res.Count + onceNum
		return res.Count, nil
	}
	next := a.currID
	a.currID++
	return next, nil
}
