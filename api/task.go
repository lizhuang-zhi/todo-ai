package api

import (
	"errors"
	"time"
	"todo-ai/common"
	"todo-ai/core"
	"todo-ai/core/logger"
	"todo-ai/model"
	"todo-ai/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

	id, err := common.TaskUUID.Get()
	if err != nil {
		logger.Errorf("TaskUUID Get error:%s", err)
		return nil, err
	}

	taskEntity := &model.Task{
		TaskID:    id,
		UserID:    request.UserID,
		Name:      request.Name,
		Type:      request.Type,
		Priority:  request.Priority,
		Date:      utils.GetTodayDate(),
		Progress:  0,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err = model.InsertTask(taskEntity)
	if err != nil {
		logger.Errorf("model.InsertTask error:%s", err)
		return nil, err
	}

	getAISuggestion(taskEntity) // 异步获取AI建议

	return id, nil
}

// go协程异步获取AI建议
func getAISuggestion(task *model.Task) {
	if task.Type != 0 { // 单日任务
		return
	}

	// TODO: 根据Date日期, 带上节假日名称等信息

	go func(t *model.Task) {
		// TODO: 异步的获取AI建议内容
		defer core.Recovery()

		// TODO: 工作流所需: name, date, 节假日名称, 优先级

		// 写入AI建议到该任务DB中
	}(task)
}

type UpdateTaskParams struct {
	TaskID   int64  `json:"task_id"`  // 任务ID
	Name     string `json:"name"`     // 任务名称
	Priority int    `json:"priority"` // 优先级(0-无 1-低 2-中 3-高)
}

// UpdateTask 创建任务
func UpdateTask(ctx *gin.Context) (interface{}, error) {
	var request UpdateTaskParams
	if err := ctx.BindJSON(&request); err != nil {
		logger.Errorf("UpdateTask BindJSON error:%s", err)
		return nil, err
	}

	if request.TaskID == 0 {
		return nil, errors.New("task_id is empty")
	}

	if request.Name == "" {
		return nil, errors.New("name is empty")
	}

	update := bson.M{"$set": bson.M{
		"name":       request.Name,
		"priority":   request.Priority,
		"updated_at": time.Now().Unix(),
	}}

	err := model.UpdateTaskByTaskID(request.TaskID, update)
	if err != nil {
		logger.Errorf("model.UpdateTaskByTaskID error:%s", err)
		return nil, err
	}

	// 获取任务信息
	task, err := model.GetTaskByTaskID(request.TaskID)
	if err != nil {
		logger.Errorf("model.GetTaskByID error:%s", err)
		return nil, err
	}

	getAISuggestion(task) // 异步获取AI建议

	return "ok", nil
}
