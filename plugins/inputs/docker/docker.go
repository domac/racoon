package docker

import (
	"github.com/phillihq/racoon/config"
)

//Docker状态信息输入配置
type DockerStatsInputConfig struct {
	config.InputConfig
	DockerURL string
}
