package ai_data

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"todo-ai/common"
	"todo-ai/core/logger"
	"todo-ai/model"
)

// 解析AI计划内容
func ParseAiPlanContent(userID int64, planCont string) error {
	operations := strings.Split(planCont, "\n")
	if len(operations) == 0 {
		return errors.New("未找到操作内容")
	}

	parentTaskID := int64(0)
	var operErr error

	for _, opra := range operations {
		if strings.Contains(opra, "[[ParentTask]]") { // 创建父任务
			parentTaskID, operErr = CreateParentTask(userID, opra)
			if operErr != nil {
				logger.Errorf("CreateParentTask失败: %v", operErr)
				return errors.New(fmt.Sprintf("CreateParentTask error:%s", operErr))
			}
			logger.Infof("parentTaskID:%d", parentTaskID)
			continue
		}

		if strings.Contains(opra, "[[SonTask]]") { // 创建子任务
			operErr := CreateSubTask(userID, opra, parentTaskID)
			if operErr != nil {
				return errors.New(fmt.Sprintf("CreateSubTask error:%s", operErr))
			}
		}
	}

	return nil
}

// 创建父任务, 格式: [[ParentTask]][add]任务名称@任务日期
func CreateParentTask(userID int64, operateLine string) (int64, error) {
	oLine := strings.TrimPrefix(operateLine, "[[ParentTask]]")
	taskInfo := strings.TrimPrefix(oLine, "[add]")
	taskInfo = strings.TrimSpace(taskInfo)
	task := strings.Split(taskInfo, "@")
	if len(task) != 2 {
		return 0, errors.New("CreateParentTask[add]任务信息格式错误")
	}

	// 添加任务
	taskName := task[0]
	taskDate := task[1]

	// 获取taskDate中的年份
	year := common.GetYearByDate(taskDate)

	parentTaskID, err := common.TaskUUID.Get()
	if err != nil {
		logger.Errorf("TaskUUID Get error:%s", err)
		return 0, err
	}

	taskEntity := &model.Task{
		TaskID:    parentTaskID,
		UserID:    userID,
		Name:      taskName,
		Type:      2, // 计划
		Priority:  0,
		Date:      taskDate,
		Year:      year,
		Progress:  0,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if err := model.InsertTask(taskEntity); err != nil {
		logger.Errorf("model.AddTask error:%s", err)
		return 0, err
	}

	return parentTaskID, nil
}

// 创建子任务, 格式: [[SonTask]][add]任务名称@任务日期|||[add]任务名称@任务日期
func CreateSubTask(userID int64, operateLine string, parentTaskID int64) error {
	if parentTaskID == 0 {
		return errors.New("parentTaskID为0")
	}

	oLine := strings.TrimPrefix(operateLine, "[[SonTask]]")
	tasks := strings.Split(oLine, "|||")
	for _, t := range tasks {
		if strings.Contains(t, "[add]") {
			taskInfo := strings.TrimPrefix(t, "[add]")
			taskInfo = strings.TrimSpace(taskInfo)
			task := strings.Split(taskInfo, "@")
			if len(task) != 2 {
				return errors.New("[add]任务信息格式错误")
			}

			// 添加任务
			taskName := task[0]
			taskDate := task[1]

			// 获取taskDate中的年份
			year := common.GetYearByDate(taskDate)

			id, err := common.TaskUUID.Get()
			if err != nil {
				logger.Errorf("TaskUUID Get error:%s", err)
				return err
			}

			taskEntity := &model.Task{
				TaskID:    id,
				UserID:    userID,
				ParentID:  parentTaskID,
				Name:      taskName,
				Type:      0, // 今日任务
				Priority:  0,
				Date:      taskDate,
				Year:      year,
				Progress:  0,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			}

			if err := model.InsertTask(taskEntity); err != nil {
				logger.Errorf("model.AddTask error:%s", err)
				return err
			}
		}
	}

	return nil
}

// 应用AI建议
