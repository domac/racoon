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
