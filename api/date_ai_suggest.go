package api

import (
	"errors"
	"time"
	"todo-ai/core/logger"
	"todo-ai/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type DateAiSuggestionRequest struct {
	UserID int64  `json:"user_id" form:"user_id"` // 用户ID
	Date   string `json:"date" form:"date"`       // 日期
}

// DateAiSuggestion 获取每日AI合理化建议
func DateAiSuggestion(ctx *gin.Context) (interface{}, error) {
	var request DateAiSuggestionRequest
	if err := ctx.Bind(&request); err != nil {
		logger.Errorf("DateAiSuggestion BindJSON error:%s", err)
		return nil, err
	}

	if request.UserID == 0 {
		return nil, errors.New("user_id is empty")
	}

	if request.Date == "" {
		return nil, errors.New("date is empty")
	}

	// 获取每日AI合理化建议
	suggest, err := model.GetDateAiSuggestByUserIDAndDate(request.UserID, request.Date)
	if err != nil {
		logger.Errorf("model.GetDateAiSuggestByUserIDAndDate error:%s", err)
		return nil, err
	}

	return suggest, nil
}

type ClickDateAiSuggestionRequest struct {
	UserID int64  `json:"user_id"` // 用户ID
	Date   string `json:"date"`    // 日期
}

// ClickDateAiSuggestion
func ClickDateAiSuggestion(ctx *gin.Context) (interface{}, error) {
	var request ClickDateAiSuggestionRequest
	if err := ctx.Bind(&request); err != nil {
		logger.Errorf("ClickDateAiSuggestion BindJSON error:%s", err)
		return nil, err
	}

	if request.UserID == 0 {
		return nil, errors.New("user_id is empty")
	}

	if request.Date == "" {
		return nil, errors.New("date is empty")
	}

	// 消除小红点
	err := model.UpsertDateAiSuggestByUserIDAndDate(request.UserID, request.Date, bson.M{"$set": bson.M{
		"show_dot":   false, // 消除小红点
		"updated_at": time.Now().Unix(),
	}})
	if err != nil {
		logger.Errorf("model.UpsertDateAiSuggestByUserIDAndDate error:%s", err)
		return nil, err
	}

	return "ok", nil
}
