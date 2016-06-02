package app

import (
	"github.com/phillihq/racoon/config"
	"github.com/phillihq/racoon/plugins/inputs/exec"
)

func init() {
	config.RegistInputHandler(exec.ModuleName, exec.InitHandler)
}
