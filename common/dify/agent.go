package dify

import (
	"context"
	"encoding/json"
	"time"
	"todo-ai/core/shttp"
)

// TODO: 待正式服更新dify版本(现在版本有Bug)
const ApiUrl = "http://10.1.7.166/v1"
const AgentSecret = "app-0K43SIwAVXhO30mwTFTsyut6"

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
			"history": history,
		},
		"query":           queryCont,
		"response_mode":   "streaming", // SSE
		"conversation_id": conversationID,
		"user":            "leo",
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

type AgentMessage struct {
	Event          string `json:"event"`           // 事件类型
	ConversationID string `json:"conversation_id"` // 会话id
	MessageID      string `json:"message_id"`
	CreatedAt      int64  `json:"created_at"`
	TaskID         string `json:"task_id"`
	ID             string `json:"id"`
	Answer         string `json:"answer"` // 回答内容
}
