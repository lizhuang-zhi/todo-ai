package dify

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"todo-ai/core/logger"
	"todo-ai/core/shttp"
)

const User = "leo"

// TODO: 待正式服更新dify版本(现在版本有Bug)
const ApiUrl = "http://10.1.7.166/v1"
const AgentSecret = "app-0K43SIwAVXhO30mwTFTsyut6" // 内网secret

// type ChatHistory struct {
// 	Messages       []*Message
// 	ConversationID string
// 	Lock           sync.Mutex
// }

// type Message struct {
// 	Role    string `json:"role"`
// 	Content string `json:"content"`
// }

// func NewChatHistory() *ChatHistory {
// 	return &ChatHistory{
// 		Messages:       []*Message{},
// 		ConversationID: "",
// 	}
// }

// func (ch *ChatHistory) AddMessage(role, content string) {
// 	ch.Lock.Lock()
// 	defer ch.Lock.Unlock()

// 	ch.Messages = append(ch.Messages, &Message{
// 		Role:    role,
// 		Content: content,
// 	})
// }

// func (ch *ChatHistory) GetHistory() string {
// 	ch.Lock.Lock()
// 	defer ch.Lock.Unlock()
// 	historyChatMessage := ""
// 	for _, msg := range ch.Messages {
// 		historyChatMessage += msg.Role + ": " + msg.Content + "\n"
// 	}
// 	return historyChatMessage
// }

// chat-message data raw
func ChatMessageDataRaw(history, queryCont, conversationID string) interface{} {
	return map[string]interface{}{
		"inputs": map[string]interface{}{
			"history":   history, // 可选(历史待办数据)
			"todayDate": time.Now().Format("2006-01-02"),
		},
		"query":           queryCont,
		"response_mode":   "streaming", // SSE
		"conversation_id": conversationID,
		"user":            User,
	}
}

// chat-message
func ChatMessage(secret string, data interface{}) ([]byte, error) {
	url := ApiUrl + "/chat-messages"

	bytes, err := json.Marshal(data)
	if err != nil {
		// logger.Errorf("[ChatMessage] json marshal with error: %v, agent[%s],", err, events.GlobalWorkflowSecretKeyMap[secret])
		return nil, err
	}

	body, err := shttp.Post(context.Background(), url, "application/json", bytes,
		shttp.WithHeader("Authorization", "Bearer "+secret), shttp.WithHTTPTimeout(3*time.Minute)).ReadAll()
	if err != nil {
		// logger.Errorf("[ChatMessage] post with error: %v, agent[%s],", err, events.GlobalWorkflowSecretKeyMap[secret])
		return nil, err
	}

	// var resp map[string]interface{}
	// if err := json.Unmarshal(body, &resp); err != nil {
	// 	// logger.Errorf("[ChatMessage] json unmarshal with error: %v, agent[%s],", err, events.GlobalWorkflowSecretKeyMap[secret])
	// 	return nil, err
	// }

	return body, nil
}

// TODO: 未使用
// 解析ChatMessage返回内容（Streaming模式)
func UnmarshalChatMessageResponse(data map[string]interface{}) string {
	event, ok := data["event"]
	if !ok {
		return ""
	}

	if event.(string) == "agent_message" { // agent类型返回
		answer, ok := data["answer"]
		if !ok {
			return ""
		}

		return answer.(string)
	}

	return ""
}

// TODO: 未使用
type AgentMessage struct {
	Event          string `json:"event"`           // 事件类型
	ConversationID string `json:"conversation_id"` // 会话id
	MessageID      string `json:"message_id"`
	CreatedAt      int64  `json:"created_at"`
	TaskID         string `json:"task_id"`
	ID             string `json:"id"`
	Answer         string `json:"answer"` // 回答内容
}

type HistoryMessageResp struct {
	Limit   int  `json:"limit"`
	HasMore bool `json:"has_more"`
	Data    []struct {
		ID             string      `json:"id"`
		ConversationID string      `json:"conversation_id"`
		Inputs         interface{} `json:"inputs"`
		Query          string      `json:"query"`
		Answer         string      `json:"answer"`
		// MessageFiles   []struct {
		// 	ID        string `json:"id"`
		// 	Type      string `json:"type"`
		// 	URL       string `json:"url"`
		// 	BelongsTo string `json:"belongs_to"`
		// } `json:"message_files"`
		// Feedback           interface{}   `json:"feedback"`
		// RetrieverResources []interface{} `json:"retriever_resources"`
		CreatedAt int64 `json:"created_at"`
		// AgentThoughts []struct {
		// 	ID           string      `json:"id"`
		// 	ChainID      interface{} `json:"chain_id"`
		// 	MessageID    string      `json:"message_id"`
		// 	Position     int         `json:"position"`
		// 	Thought      string      `json:"thought"`
		// 	Tool         string      `json:"tool"`
		// 	ToolInput    string      `json:"tool_input"`
		// 	CreatedAt    int64       `json:"created_at"`
		// 	Observation  string      `json:"observation"`
		// 	MessageFiles []string    `json:"message_files"`
		// } `json:"agent_thoughts"`
	} `json:"data"`
}

// messages 获取会话历史记录
func GetHistoryMessages(secret string, conversationID string, firstID string, limit int) (*HistoryMessageResp, error) {
	// TODO: 这里user写死了, 为leo
	url := ApiUrl + "/messages?user=" + User + "&conversation_id=" + conversationID + "&first_id=" + firstID + "&limit=" + strconv.Itoa(limit)

	body, err := shttp.Get(context.Background(), url,
		shttp.WithHeader("Authorization", "Bearer "+secret), shttp.WithHTTPTimeout(3*time.Minute)).ReadAll()
	if err != nil {
		logger.Errorf("[GetHistoryMessages] get with error: %v, conversation_id[%s], first_id[%s], limit[%d]", err, conversationID, firstID, limit)
		return nil, err
	}

	resp := &HistoryMessageResp{}
	if err := json.Unmarshal(body, resp); err != nil {
		logger.Errorf("[GetHistoryMessages] json unmarshal with error: %v, conversation_id[%s], first_id[%s], limit[%d]", err, conversationID, firstID, limit)
		return nil, err
	}

	return resp, nil
}

type HistoryConversationsResp struct {
	Limit   int  `json:"limit"`
	HasMore bool `json:"has_more"`
	Data    []struct {
		ID        string      `json:"id"`
		Name      string      `json:"name"`
		Status    string      `json:"status"`
		Inputs    interface{} `json:"inputs"`
		CreatedAt int64       `json:"created_at"`
	} `json:"data"`
}

// conversations 获取会话列表
func GetHistoryConversations(secret string, pinned bool, lastID string, limit int) (*HistoryConversationsResp, error) {
	// bool转字符串
	pinnedStr := ""
	if pinned {
		pinnedStr = "true"
	} else {
		pinnedStr = "false"
	}

	// TODO: 这里user写死了, 为leo
	url := ApiUrl + "/conversations?user=" + User + "&last_id=" + lastID + "&pinned=" + pinnedStr + "&limit=" + strconv.Itoa(limit)

	body, err := shttp.Get(context.Background(), url,
		shttp.WithHeader("Authorization", "Bearer "+secret), shttp.WithHTTPTimeout(3*time.Minute)).ReadAll()
	if err != nil {
		logger.Errorf("[GetHistoryConversations] get with error: %v, last_id[%s], pinned[%s], limit[%d]", err, lastID, pinnedStr, limit)
		return nil, err
	}

	resp := &HistoryConversationsResp{}
	if err := json.Unmarshal(body, resp); err != nil {
		logger.Errorf("[GetHistoryConversations] json unmarshal with error: %v, last_id[%s], pinned[%s], limit[%d]", err, lastID, pinnedStr, limit)
		return nil, err
	}

	return resp, nil
}
