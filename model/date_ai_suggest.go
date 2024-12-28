package model

import (
	"context"
	"todo-ai/common"
	"todo-ai/common/consts"

	"go.mongodb.org/mongo-driver/bson"
)

type DateAiSuggest struct {
	SuggestID      int64  `json:"suggest_id" bson:"suggest_id"`           // 建议ID
	UserID         int64  `json:"user_id" bson:"user_id"`                 // 用户ID
	Date           string `json:"date" bson:"date"`                       // 日期
	Status         int    `json:"status" bson:"status"`                   // 状态(0-未生成, 1-正在生成, 2-生成失败 3-生成成功)
	ShowDot        bool   `json:"show_dot" bson:"show_dot"`               // 是否显示小红点
	LastSuggestion string `json:"last_suggestion" bson:"last_suggestion"` // 最新建议
	CreatedAt      int64  `json:"created_at" bson:"created_at"`           // 创建时间
	UpdatedAt      int64  `json:"updated_at" bson:"updated_at"`           // 更新时间
}

// InsertDateAiSuggest
func InsertDateAiSuggest(dateAiSuggest *DateAiSuggest) error {
	_, err := common.Mgo.InsertOne(context.Background(), consts.CollectionDateAiSuggest, dateAiSuggest)
	return err
}

// UpdateDateAiSuggestBySuggestID
func UpdateDateAiSuggestBySuggestID(suggestID int64, update interface{}) error {
	filter := bson.M{"suggest_id": suggestID}
	_, err := common.Mgo.Update(context.Background(), consts.CollectionDateAiSuggest, filter, update)
	return err
}

// GetDateAiSuggestByUserIDAndDate
func GetDateAiSuggestByUserIDAndDate(userID int64, date string) (*DateAiSuggest, error) {
	var dateAiSuggest DateAiSuggest
	err := common.Mgo.FindOne(context.Background(), consts.CollectionDateAiSuggest, bson.M{"user_id": userID, "date": date}, &dateAiSuggest)
	if err != nil {
		return nil, err
	}
	return &dateAiSuggest, nil
}

// UpsertDateAiSuggestByUserIDAndDate
func UpsertDateAiSuggestByUserIDAndDate(userID int64, date string, update interface{}) error {
	filter := bson.M{"user_id": userID, "date": date}
	_, _, err := common.Mgo.Upsert(context.Background(), consts.CollectionDateAiSuggest, filter, update)
	return err
}
