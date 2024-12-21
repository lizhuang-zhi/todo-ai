package model

import (
	"context"
	"todo-ai/common"
	"todo-ai/common/consts"

	"go.mongodb.org/mongo-driver/bson"
)

type Task struct {
	TaskID       int64   `json:"task_id" bson:"task_id"`             // 任务ID
	UserID       int64   `json:"user_id" bson:"user_id"`             // 用户ID
	Name         string  `json:"name" bson:"name"`                   // 任务名称
	Date         string  `json:"date" bson:"date"`                   // 日期(yyyy-MM-dd)
	Priority     int     `json:"priority" bson:"priority"`           // 优先级(0-无 1-低 2-中 3-高)
	Type         int     `json:"type" bson:"type"`                   // 任务类型(0-单日任务 1-年度挑战)
	ParentID     int64   `json:"parent_id" bson:"parent_id"`         // 父任务ID
	Progress     float64 `json:"progress" bson:"progress"`           // 进度
	AiSuggestion string  `json:"ai_suggestion" bson:"ai_suggestion"` // AI建议
	CreatedAt    int64   `json:"created_at" bson:"created_at"`       // 创建时间
	UpdatedAt    int64   `json:"updated_at" bson:"updated_at"`       // 更新时间
	// Tag          string  `json:"tag" bson:"tag"`                     // 标签
}

// InsertTask
func InsertTask(task *Task) error {
	_, err := common.Mgo.InsertOne(context.Background(), consts.CollectionTask, task)
	return err
}

// UpdateTaskByTaskID
func UpdateTaskByTaskID(taskID int64, update interface{}) error {
	filter := bson.M{"task_id": taskID}
	_, err := common.Mgo.Update(context.Background(), consts.CollectionTask, filter, update)
	return err
}

// GetTaskByTaskID
func GetTaskByTaskID(taskID int64) (*Task, error) {
	var task Task
	err := common.Mgo.FindOne(context.Background(), consts.CollectionTask, bson.M{"task_id": taskID}, &task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
