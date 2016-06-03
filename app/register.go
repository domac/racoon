package app

import (
	"github.com/phillihq/racoon/config"
	"github.com/phillihq/racoon/plugins/inputs/exec"
	"github.com/phillihq/racoon/plugins/outputs/stdout"
)

func init() {
	config.RegistInputHandler(exec.ModuleName, exec.InitHandler)

	config.RegistOutputHandler(stdout.ModuleName, stdout.InitHandler)
}
