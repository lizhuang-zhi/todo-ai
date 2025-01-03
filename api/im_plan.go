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

type ChatHistoryMessagesRequest struct {
	ConversationID string `json:"conversation_id" form:"conversation_id"` // 会话ID
	FirstID        string `json:"first_id" form:"first_id"`               // 当前页的第一条聊天记录ID
	Limit          int    `json:"limit" form:"limit"`                     // 每页显示的聊天记录数
}

// ChatHistoryMessages 获取聊天历史消息
func ChatHistoryMessages(ctx *gin.Context) (interface{}, error) {
	var request ChatHistoryMessagesRequest
	if err := ctx.Bind(&request); err != nil {
		logger.Errorf("ChatHistoryMessages BindJSON error:%s", err)
		return nil, err
	}

	if request.ConversationID == "" {
		return nil, errors.New("conversation_id is empty")
	}

	if request.Limit == 0 {
		return nil, errors.New("limit is empty")
	}

	// 获取聊天历史消息
	// TODO: 替换secret
	resp, err := dify.GetHistoryMessages(dify.AgentSecret, request.ConversationID, request.FirstID, request.Limit)
	if err != nil {
		logger.Errorf("[ChatHistoryMessages] dify GetHistoryMessages with err %v, conversation_id[%s], first_id[%s], limit[%d]", err, request.ConversationID, request.FirstID, request.Limit)
		return nil, err
	}

	return resp, nil
}

type ChatHistoryConversationsRequest struct {
	Pinned bool   `json:"pinned" form:"pinned"`   // 只返回置顶 true，只返回非置顶 false
	LastID string `json:"last_id" form:"last_id"` // 当前页最后面一条记录的 ID，默认 null
	Limit  int    `json:"limit" form:"limit"`     // 一次请求返回多少条记录
}

// ChatHistoryConversations 获取会话列表
func ChatHistoryConversations(ctx *gin.Context) (interface{}, error) {
	var request ChatHistoryConversationsRequest
	if err := ctx.Bind(&request); err != nil {
		logger.Errorf("ChatHistoryConversations BindJSON error:%s", err)
		return nil, err
	}

	if request.Limit == 0 {
		return nil, errors.New("limit is empty")
	}

	// 获取聊天历史消息
	// TODO: 替换secret
	resp, err := dify.GetHistoryConversations(dify.AgentSecret, request.Pinned, request.LastID, request.Limit)
	if err != nil {
		logger.Errorf("[ChatHistoryConversations] dify GetHistoryConversations with err %v, last_id[%s], pinned[%v], limit[%d]", err, request.LastID, request.Pinned, request.Limit)
		return nil, err
	}

	return resp, nil
}
