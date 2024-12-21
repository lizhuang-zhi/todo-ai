package events

import (
	"errors"
	"strings"
	"todo-ai/common"
)

// GlobalWorkflowMap (workflow名称 -> workflow api secret)
var GlobalWorkflowMap = map[string]string{}

// GlobalWorkflowSecretKeyMap (workflow api secret -> workflow名称)
var GlobalWorkflowSecretKeyMap = map[string]string{}

func InitWorkflowCfg() error {
	workflows := common.Config.Dify.Workflow

	if len(workflows) == 0 {
		return nil
	}

	// 拆分kv
	keys, secrets, err := SplitCfg(workflows)
	if err != nil {
		return err
	}

	// 写入全局map
	for i := 0; i < len(keys); i++ {
		GlobalWorkflowMap[keys[i]] = secrets[i]
		GlobalWorkflowSecretKeyMap[secrets[i]] = keys[i]
	}

	return nil
}

func SplitCfg(workflows []string) ([]string, []string, error) {
	keys := make([]string, 0)
	secrets := make([]string, 0)

	// 拆分kv
	for _, w := range workflows {
		kv := strings.Split(w, ":")
		if len(kv) != 2 {
			return nil, nil, errors.New("workflow配置格式错误，请保证格式为key:secret")
		}

		if kv[0] == "" {
			return nil, nil, errors.New("workflow key不能为空")
		}

		if kv[1] == "" {
			return nil, nil, errors.New("workflow secret不能为空")
		}

		keys = append(keys, kv[0])
		secrets = append(secrets, kv[1])
	}

	// 判断key中是否有重复的
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] == keys[j] {
				return nil, nil, errors.New("存在重复的workflow key，请保证workflow的key唯一")
			}
		}
	}

	// 判断secret中是否有重复的
	for i := 0; i < len(secrets); i++ {
		for j := i + 1; j < len(secrets); j++ {
			if secrets[i] == secrets[j] {
				return nil, nil, errors.New("存在重复的workflow secret，请保证workflow的secret唯一")
			}
		}
	}

	return keys, secrets, nil
}
