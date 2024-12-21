package events

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"todo-ai/common"
	"todo-ai/core/logger"
	"todo-ai/core/shttp"
)

type WorkflowResp struct {
	WorkflowRunID string            `json:"workflow_run_id"`
	TaskID        string            `json:"task_id"`
	Data          *WorkflowRespData `json:"data"`
}

type WorkflowRespData struct {
	ID      string               `json:"id"`
	Outputs *WorkflowRespOutputs `json:"outputs"`
	Error   interface{}          `json:"error"`
	// 其他更多属性查看Dify 知识库API文档....
}

type WorkflowRespOutputs struct {
	Text string `json:"text"`
}

// 统一调用dify的workflow接口
func DoDifyWorkflow(secret string, data interface{}) (string, error) {
	url := common.Config.Dify.ApiUrl + "/workflows/run"

	bytes, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("workflow[%s], DoDifyWorkflow Marshal error: %v", GlobalWorkflowSecretKeyMap[secret], err)
		return "", err
	}

	body, err := shttp.Post(context.Background(), url, "application/json", bytes,
		shttp.WithHeader("Authorization", "Bearer "+secret), shttp.WithHTTPTimeout(3*time.Minute)).ReadAll()
	if err != nil {
		logger.Errorf("workflow[%s], DoDifyWorkflow Post error: %v", GlobalWorkflowSecretKeyMap[secret], err)
		return "", err
	}

	var resp WorkflowResp
	if err := json.Unmarshal(body, &resp); err != nil {
		logger.Errorf("workflow[%s], DoDifyWorkflow Unmarshal error: %v", GlobalWorkflowSecretKeyMap[secret], err)
		return "", err
	}

	if resp.Data.Error != nil {
		logger.Errorf("workflow[%s], DoDifyWorkflow error: %v", GlobalWorkflowSecretKeyMap[secret], resp.Data.Error)
		return "", errors.New("workflow error")
	}

	return resp.Data.Outputs.Text, nil
}
