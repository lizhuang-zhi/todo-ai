package api

import (
	"errors"
	"time"
	"todo-ai/common/ai_data"
	"todo-ai/core/logger"
	"todo-ai/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	if err == mongo.ErrNoDocuments { // 未找到记录
		return nil, nil
	}

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

type ApplyAiSuggestionRequest struct {
	UserID    int64  `json:"user_id"`    // 用户ID
	Date      string `json:"date"`       // 日期
	AiSuggest string `json:"ai_suggest"` // AI建议内容
}

// ApplyAiSuggestion
func ApplyAiSuggestion(ctx *gin.Context) (interface{}, error) {
	var request ApplyAiSuggestionRequest
	if err := ctx.Bind(&request); err != nil {
		logger.Errorf("ApplyAiSuggestion BindJSON error:%s", err)
		return nil, err
	}

	if request.UserID == 0 {
		return nil, errors.New("user_id is empty")
	}

	if request.Date == "" {
		return nil, errors.New("date is empty")
	}

	// TODO: 检查Ai建议的内容是修改的任务是否是该用户的任务

	// 解析AI建议内容
	result, err := ai_data.ParseAiSuggestContent(request.UserID, request.AiSuggest)
	if err != nil {
		logger.Errorf("user_id[%d], date[%s], ai_data.ParseAiSuggestContent error:%s", request.UserID, request.Date, err)
		return nil, err
	}

	for _, v := range result {
		if v != "ok" { // 存在报错
			return "执行错误, 请稍后再试~", nil
		}
	}

	// 清空AI建议内容
	err = model.UpsertDateAiSuggestByUserIDAndDate(request.UserID, request.Date, bson.M{"$set": bson.M{
		"ai_suggest": "", // 清空AI建议内容
	}})
	if err != nil {
		logger.Errorf("user_id[%d], date[%s], model.UpsertDateAiSuggestByUserIDAndDate error:%s", request.UserID, request.Date, err)
		return nil, err
	}

	return "ok", nil
}
