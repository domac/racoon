package config

type OutputContextConfig interface {
	ContextConfig
	Send(data string) (err error)
}

//输出配置
type OutputConfig struct {
	StandardConfig
}

type OutputHandler interface{}

var registedOutputHandlers = map[string]OutputHandler{}

func RegistOutputHandler(name string, handler OutputHandler) {
	registedOutputHandlers[name] = handler
}

func (self *Config) RunOutputs() (err error) {
	return self.InvokeFunc(self.runOutputs)
}

func (self *Config) runOutputs(outchan OutputCh) (err error) {
	return
}
