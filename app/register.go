package app

import (
	"github.com/phillihq/racoon/config"
	"github.com/phillihq/racoon/plugins/inputs/docker"
	"github.com/phillihq/racoon/plugins/inputs/exec"
	"github.com/phillihq/racoon/plugins/outputs/nsq"
	"github.com/phillihq/racoon/plugins/outputs/redis"
	"github.com/phillihq/racoon/plugins/outputs/stdout"
)

func init() {
	//输入注册
	config.RegistInputHandler(exec.ModuleName, exec.InitHandler)
	config.RegistInputHandler(docker.ModuleName, docker.InitHandler)

	//输出注册
	config.RegistOutputHandler(stdout.ModuleName, stdout.InitHandler)
	config.RegistOutputHandler(redis.ModuleName, redis.InitHandler)
	config.RegistOutputHandler(nsq.ModuleName, nsq.InitHandler)
}
