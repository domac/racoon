package config

type OutputContextConfig interface {
	ContextConfig
	Start()
}

//输出配置
type OutputConfig struct {
	StandardConfig
}
