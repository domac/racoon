package exec

import (
	"github.com/phillihq/racoon/config"
	"os"
)

const ModuleName = "exec"

//输入配置
type InputConfig struct {
	config.InputConfig
	Command   string   `json:"command"`
	Args      []string `json:"args,omitempty"`
	Interval  int      `json:"interval,omitempty"`
	MsgTrim   string   `json:"message_trim,omitempty"`
	MsgPrefix string   `json:"message_prefix,omitempty"`

	hostname string `json:"-"`
}

func NewInputConfig() InputConfig {
	return InputConfig{
		InputConfig: config.InputConfig{
			StandardConfig: config.StandardConfig{
				Type: ModuleName,
			},
		},
		Interval: 60,
		MsgTrim:  "\t\r\n",
	}
}

func InitHandler(configitem *config.ConfigItem) (contextConfig config.InputContextConfig, err error) {
	conf := NewInputConfig()
	if err = config.ReflectConfig(configitem, &conf); err != nil {
		return
	}
	if conf.hostname, err = os.Hostname(); err != nil {
		return
	}
	contextConfig = &conf
	return
}

func (self *InputConfig) Start() {
	println("============= exec start =============")
}
