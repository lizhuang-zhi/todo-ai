package model

import (
	"context"
	"todo-ai/common"
	"todo-ai/common/consts"
)

type Task struct {
	TaskID       int64   `json:"task_id"`       // 任务ID
	UserID       int64   `json:"user_id"`       // 用户ID
	Name         string  `json:"name"`          // 任务名称
	Date         string  `json:"date"`          // 日期(2021-01-01)
	Priority     int     `json:"priority"`      // 优先级(0-无 1-低 2-中 3-高)
	Type         int     `json:"type"`          // 任务类型(0-单日任务 1-年度挑战)
	ParentID     int64   `json:"parent_id"`     // 父任务ID
	Progress     float64 `json:"progress"`      // 进度
	AiSuggestion string  `json:"ai_suggestion"` // AI建议
	CreatedAt    int64   `json:"created_at"`    // 创建时间
	UpdatedAt    int64   `json:"updated_at"`    // 更新时间
	// Tag          string  `json:"tag"`           // 标签(暂时不用)
}

// InsertTask
func InsertTask(task *Task) error {
	_, err := common.Mgo.InsertOne(context.Background(), consts.CollectionTask, task)
	return err
}
