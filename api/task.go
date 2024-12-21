package api

import (
	"errors"
	"time"
	"todo-ai/core"
	"todo-ai/core/logger"
	"todo-ai/model"
	"todo-ai/utils"

	"github.com/gin-gonic/gin"
)

type CreateTaskParams struct {
	UserID   int64  `json:"user_id"`  // 用户ID(暂时不用, 默认传1)
	Name     string `json:"name"`     // 任务名称
	Type     int    `json:"type"`     // 任务类型(0-单日任务 1-年度挑战)
	Priority int    `json:"priority"` // 优先级(0-无 1-低 2-中 3-高)
}

// CreateTask 创建任务
func CreateTask(ctx *gin.Context) (interface{}, error) {
	var request CreateTaskParams
	if err := ctx.BindJSON(&request); err != nil {
		logger.Errorf("CreateTask BindJSON error:%s", err)
		return nil, err
	}

	if request.UserID == 0 {
		return nil, errors.New("user_id is empty")
	}

	if request.Name == "" {
		return nil, errors.New("name is empty")
	}

	taskEntity := &model.Task{
		UserID:    request.UserID,
		Name:      request.Name,
		Type:      request.Type,
		Priority:  request.Priority,
		Date:      utils.GetTodayDate(),
		Progress:  0,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err := model.InsertTask(taskEntity)
	if err != nil {
		logger.Errorf("model.DeleteRecordByUserIDAndType error:%s", err)
		return nil, err
	}

	getAISuggestion(taskEntity) // 异步获取AI建议

	return "ok", nil
}

// go协程异步获取AI建议
func getAISuggestion(task *model.Task) {
	if task.Type != 0 { // 单日任务
		return
	}

	go func(t *model.Task) {
		// TODO: 异步的获取AI建议内容
		defer core.Recovery()

		// 写入AI建议到该任务DB中
	}(task)
}
