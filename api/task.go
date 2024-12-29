package api

import (
	"errors"
	"time"
	"todo-ai/common"
	"todo-ai/common/ai_data"
	"todo-ai/core"
	"todo-ai/core/logger"
	"todo-ai/events"
	"todo-ai/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type CreateTaskParams struct {
	UserID   int64  `json:"user_id"`  // 用户ID(暂时不用, 默认传1)
	Name     string `json:"name"`     // 任务名称
	Type     int    `json:"type"`     // 任务类型(0-单日任务 1-年度挑战)
	Priority int    `json:"priority"` // 优先级(0-无 1-低 2-中 3-高)
	Date     string `json:"date"`     // 日期(yyyy-MM-dd)
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

	if request.Date == "" {
		return nil, errors.New("date is empty")
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
		Date:      request.Date,
		Year:      time.Now().Year(),
		Progress:  0,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err = model.InsertTask(taskEntity)
	if err != nil {
		logger.Errorf("model.InsertTask error:%s", err)
		return nil, err
	}

	getAISuggestion(taskEntity)     // 异步获取AI建议
	getDateAISuggestion(taskEntity) // 异步获取AI今日规划合理化建议

	return id, nil
}

// 异步获取AI今日规划合理化建议
func getDateAISuggestion(task *model.Task) {
	if task.Type != 0 { // 单日任务
		return
	}

	err := model.UpsertDateAiSuggestByUserIDAndDate(task.UserID, task.Date,
		bson.M{"$set": bson.M{
			"last_suggestion": "",
			"show_dot":        false, // 不显示小红点
			"status":          1,     // 生成中,
			"updated_at":      time.Now().Unix(),
		}})
	if err != nil {
		logger.Errorf("model.UpsertDateAiSuggestByUserIDAndDate error:%s", err)
	}

	go func(t *model.Task) {
		// 异步的获取AI建议内容
		defer core.Recovery()

		secret, ok := events.GlobalWorkflowMap["[todo]当日任务排期合理化工作流"]
		if !ok {
			logger.Errorf("GlobalWorkflowMap secret not found: %s", "[todo]当日任务排期合理化工作流")
			return
		}

		// 获取今日代办任务, 并整理格式
		todayTasksStr, err := ai_data.GetTodayTasksData(t.UserID, t.Date)
		if err != nil {
			logger.Errorf("user_id[%s], date[%s], GetTodayTasksData error:%s", t.UserID, t.Date, err)
			return
		}

		// 近10天的历史待办事项
		historyStr, err := ai_data.GetHistoryTasksData(t.UserID)
		if err != nil {
			logger.Errorf("user_id[%s], date[%s], GetTodayTasksData error:%s", t.UserID, t.Date, err)
			return
		}

		// 今天的日期、节日等情况
		holidayStr, err := ai_data.GetHolidayData(t.Date)
		if err != nil {
			logger.Errorf("user_id[%s], date[%s], GetTodayTasksData error:%s", t.UserID, t.Date, err)
			return
		}

		dateAiSuggest, err := DoDifyWorkflowDateAiSuggest(secret, todayTasksStr, historyStr, holidayStr)
		if err != nil {
			logger.Errorf("DoDifyWorkflow error:%s", err)
			return
		}

		err = model.UpsertDateAiSuggestByUserIDAndDate(task.UserID, task.Date,
			bson.M{"$set": bson.M{
				"last_suggestion": dateAiSuggest,
				"status":          3,    // 生成成功,
				"show_dot":        true, // 显示小红点
				"updated_at":      time.Now().Unix(),
			}})
		if err != nil {
			logger.Errorf("model.UpsertDateAiSuggestByUserIDAndDate error:%s", err)
		}
	}(task)
}

func DoDifyWorkflowDateAiSuggest(secret string, todayTasksStr, historyStr, holidayStr interface{}) (string, error) {
	data := map[string]interface{}{
		"inputs": map[string]interface{}{
			"todayTasks": todayTasksStr,
			"history":    historyStr,
			"day":        holidayStr,
		}, // 其他参数
		"response_mode": "blocking",       // blocking 阻塞、non_blocking 非阻塞
		"user":          "todo-ai-server", // 必须要填写
	}

	return events.DoDifyWorkflow(secret, data)
}

// go协程异步获取AI建议
func getAISuggestion(task *model.Task) {
	if task.Type != 0 { // 单日任务
		return
	}

	// 根据Date日期, 带上节假日名称等信息(暂时不带)

	go func(t *model.Task) {
		// 异步的获取AI建议内容
		defer core.Recovery()

		secret, ok := events.GlobalWorkflowMap["[todo]当日任务AI建议工作流"]
		if !ok {
			logger.Errorf("GlobalWorkflowMap secret not found: %s", "[todo]当日任务AI建议工作流")
			return
		}

		aiCont, err := DoDifyWorkflow(secret, task.Name)
		if err != nil {
			logger.Errorf("DoDifyWorkflow error:%s", err)
			return
		}

		err = model.UpdateTaskByTaskID(t.TaskID, bson.M{"$set": bson.M{"ai_suggestion": aiCont}})
		if err != nil {
			logger.Errorf("model.UpdateTaskByTaskID error:%s", err)
		}
	}(task)
}

// 请求线上工作流获取分析后的AI结果
func DoDifyWorkflow(secret string, cont interface{}) (string, error) {
	data := map[string]interface{}{
		"inputs": map[string]interface{}{
			"todoItem": cont,
		}, // 其他参数
		"response_mode": "blocking",       // blocking 阻塞、non_blocking 非阻塞
		"user":          "todo-ai-server", // 必须要填写
	}

	return events.DoDifyWorkflow(secret, data)
}

type UpdateTaskParams struct {
	TaskID   int64  `json:"task_id"`  // 任务ID
	Name     string `json:"name"`     // 任务名称
	Priority int    `json:"priority"` // 优先级(0-无 1-低 2-中 3-高)
	// Date     string `json:"date"`     // 日期(yyyy-MM-dd)  // TODO: 改到某日期
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

	// TODO: 改到某日期, 还要同步修改Date和Year字段

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

	getAISuggestion(task)     // 异步获取AI建议
	getDateAISuggestion(task) // 异步获取AI今日规划合理化建议

	return "ok", nil
}

type DeleteTaskParams struct {
	TaskID int64 `json:"task_id"` // 任务ID
}

// DeleteTask 删除任务
func DeleteTask(ctx *gin.Context) (interface{}, error) {
	var request DeleteTaskParams
	if err := ctx.BindJSON(&request); err != nil {
		logger.Errorf("DeleteTask BindJSON error:%s", err)
		return nil, err
	}

	if request.TaskID == 0 {
		return nil, errors.New("task_id is empty")
	}

	err := model.DeleteTaskByTaskID(request.TaskID)
	if err != nil {
		logger.Errorf("model.DeleteTaskByTaskID error:%s", err)
		return nil, err
	}

	return "ok", nil
}

type FinishedTaskarams struct {
	TaskID int64 `json:"task_id"` // 任务ID
}

// FinishedTask 完成任务
func FinishedTask(ctx *gin.Context) (interface{}, error) {
	var request FinishedTaskarams
	if err := ctx.BindJSON(&request); err != nil {
		logger.Errorf("FinishedTask BindJSON error:%s", err)
		return nil, err
	}

	if request.TaskID == 0 {
		return nil, errors.New("task_id is empty")
	}

	// 获取任务信息
	task, err := model.GetTaskByTaskID(request.TaskID)
	if err != nil {
		logger.Errorf("model.GetTaskByID error:%s", err)
		return nil, err
	}

	update := bson.M{"$set": bson.M{
		"progress":   1,
		"updated_at": time.Now().Unix(),
	}}

	err = model.UpdateTaskByTaskID(request.TaskID, update)
	if err != nil {
		logger.Errorf("model.UpdateTaskByTaskID error:%s", err)
		return nil, err
	}

	// 单日任务
	if task.Type == 0 && task.ParentID != 0 {
		calParentTaskProgress(task.ParentID)
	} else if task.Type == 1 { // 年度挑战任务
		calSubTaskProgress(task.TaskID)
	}

	return "ok", nil
}

// 重新计算父任务的progress
func calParentTaskProgress(parentID int64) {
	// 获取所有子任务
	subTasks, err := model.GetSubTaskByParentID(parentID)
	if err != nil {
		logger.Errorf("model.GetSubTaskByParentID error:%s", err)
	}

	// 计算progress
	finishedCnt := 0
	for _, task := range subTasks {
		if task.Progress == 1 {
			finishedCnt++
		}
	}

	progress := float64(finishedCnt) / float64(len(subTasks))

	// 更新父任务的progress
	update := bson.M{"$set": bson.M{
		"progress":   progress,
		"updated_at": time.Now().Unix(),
	}}
	err = model.UpdateTaskByTaskID(parentID, update)
	if err != nil {
		logger.Errorf("model.UpdateTaskByTaskID error:%s", err)
	}
}

// 所有关联子任务的progress设置为1(完成)
func calSubTaskProgress(taskID int64) {
	// 获取所有子任务
	subTasks, err := model.GetSubTaskByParentID(taskID)
	if err != nil {
		logger.Errorf("model.GetSubTaskByParentID error:%s", err)
	}

	// 更新子任务的progress
	for _, task := range subTasks {
		update := bson.M{"$set": bson.M{
			"progress":   1,
			"updated_at": time.Now().Unix(),
		}}
		err = model.UpdateTaskByTaskID(task.TaskID, update)
		if err != nil {
			logger.Errorf("model.UpdateTaskByTaskID error:%s", err)
		}
	}
}

type ListTaskRequest struct {
	UserID int64  `json:"user_id" form:"user_id"` // 用户ID
	Type   int    `json:"type" form:"type"`       // 任务类型(0-单日任务 1-年度挑战)
	Date   string `json:"date" form:"date"`       // 日期(yyyy-MM-dd)
	Year   int    `json:"year" form:"year"`       // 年份
}

// ListTask 获取任务列表
func ListTask(ctx *gin.Context) (interface{}, error) {
	var request ListTaskRequest
	if err := ctx.Bind(&request); err != nil {
		logger.Errorf("ListTask BindJSON error:%s", err)
		return nil, err
	}

	if request.UserID == 0 {
		return nil, errors.New("user_id is empty")
	}

	if request.Date == "" {
		return nil, errors.New("date is empty")
	}

	if request.Year == 0 {
		return nil, errors.New("year is empty")
	}

	if request.Type == 0 { // 单日任务
		return listDayTask(request)
	} else if request.Type == 1 { // 年度挑战任务
		return listYearTask(request)
	}

	return nil, nil
}

// 获取单日任务列表
func listDayTask(request ListTaskRequest) (interface{}, error) {
	// 获取当日任务
	tasks, err := model.GetTaskByUserIDAndDate(request.UserID, request.Date)
	if err != nil {
		logger.Errorf("model.GetTaskByUserIDAndDate error:%s", err)
		return nil, err
	}

	return tasks, nil
}

// 获取年度挑战任务列表
func listYearTask(request ListTaskRequest) (interface{}, error) {
	// 获取年度挑战任务
	tasks, err := model.GetTaskByUserIDAndYear(request.UserID, request.Year)
	if err != nil {
		logger.Errorf("model.GetTaskByUserIDAndYear error:%s", err)
		return nil, err
	}
	return tasks, nil
}

type GetTaskDetailRequest struct {
	TaskID int64 `json:"task_id" form:"task_id"` // 任务ID
}

// GetTaskDetail 根据任务id获取任务详情
func GetTaskDetail(ctx *gin.Context) (interface{}, error) {
	var request GetTaskDetailRequest
	if err := ctx.Bind(&request); err != nil {
		logger.Errorf("GetTaskDetail BindJSON error:%s", err)
		return nil, err
	}

	if request.TaskID == 0 {
		return nil, errors.New("task_id is empty")
	}

	task, err := model.GetTaskByTaskID(request.TaskID)
	if err != nil {
		logger.Errorf("task_id[%d], model.GetTaskByTaskID error:%s", request.TaskID, err)
		return nil, err
	}

	return task, nil
}
