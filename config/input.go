package config

//输入配置接口
type InputContextConfig interface {
	ContextConfig
	Start()
}

type InputConfig struct {
	StandardConfig
}

type InputHandler interface{}

var AllInputHandlers = map[string]InputHandler{}

func (self *Config) RunInputs() error {
	return self.Invoke(self.runInputs)
}

func (self *Config) runInputs(inChan InputCh) error {

	println("------------------runInputs")

	return nil
}
