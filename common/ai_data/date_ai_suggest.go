package ai_data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"todo-ai/common"
	"todo-ai/common/holiday"
	"todo-ai/core/logger"
	"todo-ai/model"

	"go.mongodb.org/mongo-driver/bson"
)

// 获取今日待办任务数据(ai所需格式化)
func GetTodayTasksData(userID int64, date string) (string, error) {
	// 获取当日任务
	tasks, err := model.GetTaskByUserIDAndDate(userID, date)
	if err != nil {
		logger.Errorf("model.GetTaskByUserIDAndDate error:%s", err)
		return "", err
	}

	data := "今日代办任务如下：\n"

	// 格式化数据
	for _, task := range tasks {
		if task.Type != 0 { // 非今日任务
			continue
		}

		isFinish := ""
		if task.Progress == 1 {
			isFinish = "已完成"
		} else {
			isFinish = "未完成"
		}

		data += fmt.Sprintf("任务id: %d,任务名称：%s,任务状态：%s\n", task.TaskID, task.Name, isFinish)
	}

	return data, nil
}

// 近10天的历史待办事项
func GetHistoryTasksData(userID int64) (string, error) {
	// 获取近10天的历史任务
	date := time.Now().AddDate(0, 0, -10).Format("2006-01-02")
	recentlyTasks, err := model.GetTaskByUserIDAndDateDesc(userID, date)
	if err != nil {
		logger.Errorf("model.GetTaskByUserIDAndRecentlyDate error:%s", err)
		return "", err
	}

	data := "近10天的历史待办事项如下：\n"

	// 格式化数据
	for _, task := range recentlyTasks {
		if task.Type != 0 { // 非今日任务
			continue
		}

		isFinish := ""
		if task.Progress == 1 {
			isFinish = "已完成"
		} else {
			isFinish = "未完成"
		}

		data += fmt.Sprintf("任务id: %d,日期: %s,任务名称：%s,任务状态：%s\n", task.TaskID, task.Date, task.Name, isFinish)
	}

	return data, nil
}

// 获取日期对应的节假日信息
func GetHolidayData(date string) (string, error) {
	// 判断是否是节假日
	for d, v := range holiday.LegalHolidays {
		if d == date {
			return fmt.Sprintf("%s是法定节假日, 节日名称：%s", date, v), nil
		}
	}

	// 判断date是否是工作日
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", errors.New("日期格式错误")
	}

	resDayData := ""
	if holiday.IsWorkday(dateTime) {
		resDayData += fmt.Sprintf("%s是工作日", date)
	} else {
		resDayData += fmt.Sprintf("%s是周末", date)
	}

	// 判断是否是普通节日
	for d, v := range holiday.NormalHolidays {
		if d == date {
			resDayData += fmt.Sprintf(",同时%s是普通节日, 节日名称：%s", date, v)
			break
		}
	}

	return resDayData, nil
}

// 解析AI建议内容
func ParseAiSuggestContent(userID int64, suggest string) ([]string, error) {
	// 找到```包裹内容```中的包裹内容
	startIndex := strings.Index(suggest, "```")
	if startIndex == -1 {
		return nil, errors.New("未找到开头的```")
	}

	endIndex := strings.LastIndex(suggest, "```")
	if endIndex == -1 {
		return nil, errors.New("未找到结尾的```")
	}

	// 获取```包裹内容```
	content := suggest[startIndex+3 : endIndex]

	operations := strings.Split(content, "\n")
	if len(operations) == 0 {
		return nil, errors.New("未找到操作内容")
	}

	operateResult := make([]string, 0)

	for _, opra := range operations {
		if strings.Contains(opra, "[[SplitTask]]") {
			operErr := SplitTask(userID, opra)
			if operErr != nil {
				errStr := fmt.Sprintf("SplitTask error:%s", operErr)
				operateResult = append(operateResult, errStr)
			} else {
				operateResult = append(operateResult, "ok")
			}
		} else if strings.Contains(opra, "[[UpdateNameTask]]") {
			operErr := UpdateNameTask(opra)
			if operErr != nil {
				errStr := fmt.Sprintf("UpdateNameTask error:%s", operErr)
				operateResult = append(operateResult, errStr)
			} else {
				operateResult = append(operateResult, "ok")
			}
		} else if strings.Contains(opra, "[[UpdateDateTask]]") {
			operErr := UpdateDateTask(opra)
			if operErr != nil {
				errStr := fmt.Sprintf("UpdateDateTask error:%s", operErr)
				operateResult = append(operateResult, errStr)
			} else {
				operateResult = append(operateResult, "ok")
			}
		}
	}

	return operateResult, nil
}

// 拆分任务, 格式: [[SplitTask]][delete]任务id|||[add]任务名称@任务日期|||[add]任务名称@任务日期
func SplitTask(userID int64, operateLine string) error {
	oLine := strings.TrimPrefix(operateLine, "[[SplitTask]]")
	tasks := strings.Split(oLine, "|||")
	for _, t := range tasks {
		if strings.Contains(t, "[delete]") {
			taskID := strings.TrimPrefix(t, "[delete]")
			taskID = strings.TrimSpace(taskID)

			taskIDInt64, err := strconv.ParseInt(taskID, 10, 64)
			if err != nil {
				return err
			}

			// 删除任务
			if err := model.DeleteTaskByTaskID(taskIDInt64); err != nil {
				logger.Errorf("model.DeleteTaskByTaskID error:%s", err)
				return err
			}
		} else if strings.Contains(t, "[add]") {
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

// 更新任务名称, 格式: [[UpdateNameTask]][update_name]任务id@新任务名称
func UpdateNameTask(operateLine string) error {
	oLine := strings.TrimPrefix(operateLine, "[[UpdateNameTask]]")
	tasks := strings.Split(oLine, "|||")
	for _, t := range tasks {
		if strings.Contains(t, "[update_name]") {
			taskInfo := strings.TrimPrefix(t, "[update_name]")
			taskInfo = strings.TrimSpace(taskInfo)
			task := strings.Split(taskInfo, "@")
			if len(task) != 2 {
				return errors.New("[update_name]任务信息格式错误")
			}

			// 更新任务名称
			taskID := task[0]
			newTaskName := task[1]

			taskIDInt64, err := strconv.ParseInt(taskID, 10, 64)
			if err != nil {
				return err
			}

			// 更新任务名称
			if err := model.UpdateTaskByTaskID(taskIDInt64, bson.M{"$set": bson.M{"name": newTaskName}}); err != nil {
				logger.Errorf("model.UpdateTaskByTaskID error:%s", err)
				return err
			}
		}
	}

	return nil
}

// 更新任务日期, 格式: [[UpdateDateTask]][update_date]任务id@新日期
func UpdateDateTask(operateLine string) error {
	oLine := strings.TrimPrefix(operateLine, "[[UpdateDateTask]]")
	tasks := strings.Split(oLine, "|||")
	for _, t := range tasks {
		if strings.Contains(t, "[update_date]") {
			taskInfo := strings.TrimPrefix(t, "[update_date]")
			taskInfo = strings.TrimSpace(taskInfo)
			task := strings.Split(taskInfo, "@")
			if len(task) != 2 {
				return errors.New("[update_date]任务信息格式错误")
			}

			// 更新任务日期
			taskID := task[0]
			newTaskDate := task[1]

			taskIDInt64, err := strconv.ParseInt(taskID, 10, 64)
			if err != nil {
				return err
			}

			// 更新任务日期
			if err := model.UpdateTaskByTaskID(taskIDInt64, bson.M{"$set": bson.M{"date": newTaskDate}}); err != nil {
				logger.Errorf("model.UpdateTaskByTaskID error:%s", err)
				return err
			}
		}
	}

	return nil
}

// 应用AI建议
