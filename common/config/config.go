package config

import "strings"

type Server struct {
	System struct {
		ServerType    int64  `mapstructure:"server_type"` // 服务器类型 0表示国际服 1国服
		Server        string `mapstructure:"server"`      // 服务器名称  未设置则使用localip
		TCPPort       string `mapstructure:"tcp_port"`
		AuthOpen      bool   `mapstructure:"auth_open" usage:"auth_open"` // 权限开关
		AuthAccessKey string `mapstructure:"auth_access_key" usage:"auth_access_key"`
		OpenAPIToken  string `mapstructure:"openapi_token" usage:"openapi_token"`     // OpenAPI Token
		TmpDir        string `mapstructure:"tmp_dir" usage:"tmp_dir"`                 // 临时数据存放目录
		Cloud         string `mapstructure:"cloud" usage:"deploy cloud"`              // 部署所在云厂商
		Shadow        bool   `mapstructure:"shadow" usage:"shadow"`                   // 是否是影子环境
		MobileBaseURL string `mapstructure:"mobile_base_url" usage:"mobile_base_url"` // mobile基础地址
	}

	// 事件分发器
	Events struct {
		Cron bool `mapstructure:"cron"` // 是否开启定时任务
	} `mapstructure:"events"`

	// MongoDB
	MongoDB struct {
		Instance string `mapstructure:"instance"` // MongoDB URL
	} `mapstructure:"mongodb"`

	// Redis
	Redis struct {
		Instance  string `mapstructure:"instance"` //
		IsCluster bool   `mapstructure:"isCluster"`
		Prefix    string `mapstructure:"prefix"` //
	} `mapstructure:"redis"`

	// 日志配置
	Logging struct {
		Level      string `mapstructure:"level"`       // 日志等级
		Color      bool   `mapstructure:"color"`       // 是否在终端显示颜色
		Path       string `mapstructure:"path"`        // 日志路径
		Encoding   string `mapstructure:"encoding"`    // 日志编码(console, json)
		MaxSize    int    `mapstructure:"max_size"`    // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups int    `mapstructure:"max_backups"` // 保留旧文件的最大个数
		MaxAge     int    `mapstructure:"max_age"`     // 保留旧文件的最大天数
		Buffer     bool   `mapstructure:"buffer"`      // 缓写模式
	} `mapstructure:"logging"`

	// Dify配置
	Dify struct {
		ApiUrl   string   `mapstructure:"api_url"`  // Dify API地址
		Workflow []string `mapstructure:"workflow"` // Dify工作流KV
	} `mapstructure:"dify"`
}

// GetNameForURL 获取用于url的名称
func (s *Server) GetNameForURL() string {
	name := strings.ReplaceAll(s.System.Server, "_", "-") // ios的url不支持下划线
	name = strings.ReplaceAll(name, " ", "-")             // ios的url不支持空格
	return name
}
