package api

import (
	"errors"
	"todo-ai/core/logger"

	"github.com/gin-gonic/gin"
)

type DayGenerateAiPlanParams struct {
	Tag       string `json:"tag"`        // 标签
	TotalTime int    `json:"total_time"` // 总时间(单位:秒)
}

// DayGenerateAiPlan 生成AI当日规划
func DayGenerateAiPlan(ctx *gin.Context) (interface{}, error) {
	var request DayGenerateAiPlanParams
	if err := ctx.BindJSON(&request); err != nil {
		logger.Errorf("DayGenerateAiPlan BindJSON error:%s", err)
		return nil, err
	}

	if request.Tag == "" {
		return nil, errors.New("tag is empty")
	}

	if request.TotalTime == 0 {
		return nil, errors.New("total_time is empty")
	}

	// TODO: 携带 Date、Tag、TotalTime + 10天完成情况的历史数据
	aiRes := ""

	// TODO: 处理展示格式

	return aiRes, nil
}

type DayApplyAiPlanParams struct {
	AiGenCont string `json:"ai_gen_cont"` // AI生成内容
	Date      string `json:"date"`        // 日期
}

// DayApplyAiPlan 应用AI当日规划
func DayApplyAiPlan(ctx *gin.Context) (interface{}, error) {
	var request DayApplyAiPlanParams
	if err := ctx.BindJSON(&request); err != nil {
		logger.Errorf("DayApplyAiPlan BindJSON error:%s", err)
		return nil, err
	}

	if request.AiGenCont == "" {
		return nil, errors.New("ai_gen_cont is empty")
	}

	if request.Date == "" {
		return nil, errors.New("date is empty")
	}

	// TODO: 解析AiGenCont

	// TODO: 遍历创建当日任务, 并写入DB

	return "ok", nil
}
