package dify

import (
	"encoding/json"
	"strings"
	"testing"
	"todo-ai/common"
	"todo-ai/common/config"
	"todo-ai/events"
)

func TestChatMessage(t *testing.T) {
	common.Config = &config.Server{}
	common.Config.Dify.ApiUrl = "http://dify.fpsops.com/v1"
	common.Config.Dify.Workflow = []string{"对话Agent:app-QFo9nssgEuLDiL6sSK5hRWO2"}

	err := events.InitWorkflowCfg()
	if err != nil {
		t.Error(err)
		return
	}

	secret, ok := events.GlobalWorkflowMap["对话Agent"]
	if !ok {
		t.Error("secret not found")
		return
	}

	t.Log("secret:", secret)

	// 构建历史对话
	historyChatMessage := NewChatHistory()
	historyChatMessage.AddMessage("user", "南极旅行计划")
	historyChatMessage.AddMessage("agent", "好的，收到，那么您的计划是什么时候出发呢？")
	historyChatMessage.AddMessage("user", "2025-01-01")
	historyChatMessage.AddMessage("agent", "好的，收到，那么您的计划存在间隔吗？")

	data := ChatMessageDataRaw(historyChatMessage.GetHistory(), "间隔1天", "0e65261c-8195-4276-82b3-9d14c17c13f1")

	res, err := ChatMessage(secret, data)
	if err != nil {
		t.Error(err)
		return
	}

	result := string(res)

	// 按照\n\n分割
	resultChunk := strings.Split(result, "\n\n")

	wholeAnswer := ""

	for _, chunk := range resultChunk {
		// 去除data:开头
		chunk = strings.TrimPrefix(chunk, "data:")
		// 去除空格
		chunk = strings.TrimSpace(chunk)
		t.Log("Chunk: ", chunk)

		if chunk == "" {
			continue
		}

		// 字符串转[]byte
		chunkBytes := []byte(chunk)
		var agentMsg map[string]interface{}
		err := json.Unmarshal(chunkBytes, &agentMsg)
		if err != nil {
			t.Errorf("json unmarshal error: %v", err)
			return
		}

		event, ok := agentMsg["event"]
		if !ok {
			t.Errorf("event not found, %v", agentMsg)
			continue
		}

		if event.(string) == "agent_message" {
			a, ok := agentMsg["answer"].(string)
			if !ok {
				t.Errorf("answer not found, %v", agentMsg)
				continue
			}

			wholeAnswer += a
		}
	}

	t.Log("Whole Answer: ", wholeAnswer)
}
