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
	Year         int     `json:"year" bson:"year"`                   // 年份
	Priority     int     `json:"priority" bson:"priority"`           // 优先级(0-无 1-低 2-中 3-高)
	Type         int     `json:"type" bson:"type"`                   // 任务类型(0-单日任务 1-年度挑战 2-计划)
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

// DeleteTaskByTaskID
func DeleteTaskByTaskID(taskID int64) error {
	_, err := common.Mgo.Delete(context.Background(), consts.CollectionTask, bson.M{"task_id": taskID})
	return err
}

// GetSubTaskByParentID
func GetSubTaskByParentID(parentID int64) ([]*Task, error) {
	var tasks []*Task
	err := common.Mgo.Find(context.Background(), consts.CollectionTask, bson.M{"parent_id": parentID}, nil, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTaskByUserIDAndDate  根据用户ID和日期获取当日任务
func GetTaskByUserIDAndDate(userID int64, date string) ([]*Task, error) {
	var tasks []*Task
	err := common.Mgo.Find(context.Background(), consts.CollectionTask, bson.M{"user_id": userID, "date": date, "type": 0}, nil, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTaskByUserIDAndYear 根据用户ID和年份获取年度挑战任务
func GetTaskByUserIDAndYear(userID int64, year int) ([]*Task, error) {
	var tasks []*Task
	err := common.Mgo.Find(context.Background(), consts.CollectionTask, bson.M{"user_id": userID, "year": year, "type": 1}, nil, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// 根据用户ID获取总任务数量
func GetTotalTaskLenByUserID(userID int64) (int64, error) {
	total, err := common.Mgo.Count(context.Background(), consts.CollectionTask, bson.M{"user_id": userID})
	if err != nil {
		return 0, err
	}
	return total, nil
}

// 根据用户ID获取已完成任务数量
func GetFinishedTaskLenByUserID(userID int64) (int64, error) {
	total, err := common.Mgo.Count(context.Background(), consts.CollectionTask, bson.M{"user_id": userID, "progress": 1})
	if err != nil {
		return 0, err
	}
	return total, nil
}

// 根据用户ID和30天前的日期, 获取用户近30天的任务数据
func GetTaskByUserIDAndRecentlyDate(userID int64, date string) ([]*Task, error) {
	var tasks []*Task
	err := common.Mgo.Find(context.Background(), consts.CollectionTask, bson.M{"user_id": userID, "date": bson.M{"$gte": date}}, nil, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// 根据用户ID和日期获取任务, 按照日期倒序排序
// TODO: 这里有bug, 不是取近n天, 而是历史数据的80条数据
func GetTaskByUserIDAndDateDesc(userID int64, date string) ([]*Task, error) {
	var tasks []*Task
	err := common.Mgo.FindSortDesRange(context.Background(), consts.CollectionTask, bson.M{"user_id": userID, "date": bson.M{"$gte": date}}, "date", 0, 80, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
