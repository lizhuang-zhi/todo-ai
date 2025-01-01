package api

import (
	"errors"
	"todo-ai/common/ai_data"
	"todo-ai/common/dify"
	"todo-ai/core/logger"

	"github.com/gin-gonic/gin"
)

type ChatMessageParams struct {
	UserID         int64  `json:"user_id"`         // 用户ID(暂时不用, 默认传1)
	Query          string `json:"query"`           // 用户输入
	ConversationID string `json:"conversation_id"` // 会话ID(创建时不需要传, 之后的消息都需要传)
}

// ChatMessage
func ChatMessage(ctx *gin.Context) (interface{}, error) {
	var request ChatMessageParams
	if err := ctx.BindJSON(&request); err != nil {
		logger.Errorf("ChatMessage BindJSON error:%s", err)
		return nil, err
	}

	if request.UserID == 0 {
		return nil, errors.New("user_id is empty")
	}

	if request.Query == "" {
		return nil, errors.New("query is empty")
	}

	history := "" // 历史完成任务情况
	// 创建新的聊天任务
	data := dify.ChatMessageDataRaw(history, request.Query, request.ConversationID)

	// TODO: 替换secret
	btyes, err := dify.ChatMessage(dify.AgentSecret, data)
	if err != nil {
		return nil, err
	}

	return string(btyes), nil
}

type DayApplyAiPlanParams struct {
	UserID        int64  `json:"user_id"`     // 用户ID
	AiPlanGenCont string `json:"ai_gen_cont"` // AI生成的计划内容
}

// DayApplyAiPlan 应用AI当日规划
func DayApplyAiPlan(ctx *gin.Context) (interface{}, error) {
	var request DayApplyAiPlanParams
	if err := ctx.BindJSON(&request); err != nil {
		logger.Errorf("DayApplyAiPlan BindJSON error:%s", err)
		return nil, err
	}

	if request.AiPlanGenCont == "" {
		return nil, errors.New("ai_plan_gen_cont is empty")
	}

	// 解析AI计划
	err := ai_data.ParseAiPlanContent(request.UserID, request.AiPlanGenCont)
	if err != nil {
		logger.Errorf("[DayApplyAiPlan] ParseAiSuggestContent error:%s, user_id:%d, ai_gen_cont:%s", err, request.UserID, request.AiPlanGenCont)
		return nil, err
	}

	return "ok", nil
}
