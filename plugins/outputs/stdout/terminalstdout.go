package stdout

import (
	"fmt"
	"github.com/phillihq/racoon/config"
	"github.com/phillihq/racoon/config/gather"
)

const (
	ModuleName = "stdout"
)

type OutputConfig struct {
	config.OutputConfig
}

func NewOutputConfig() OutputConfig {
	return OutputConfig{
		OutputConfig: config.OutputConfig{
			StandardConfig: config.StandardConfig{
				Type: ModuleName,
			},
		},
	}
}

func InitHandler(configitem *config.ConfigItem) (contextConfig config.OutputContextConfig, err error) {
	conf := NewOutputConfig()
	if err = config.ReflectConfig(configitem, &conf); err != nil {
		return
	}
	contextConfig = &conf
	return
}

func (self *OutputConfig) Send(data gather.GatherData) (err error) {
	fdata, err := data.MarshalIndent()
	if err != nil {
		return
	}
	fmt.Println(string(fdata))
	return
}
