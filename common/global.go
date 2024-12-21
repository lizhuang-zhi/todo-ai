package common

import (
	"todo-ai/common/config"
	"todo-ai/core"
)

var (
	Config *config.Server

	Mgo *core.Mongodb

	UserUUID core.IntUUID // 用户信息id
	TaskUUID core.IntUUID // 记录id
)
