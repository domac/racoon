package config

type FilterContextConfig interface {
	ContextConfig
	Start()
}

//过滤器配置信息
type FilterConfig struct {
	StandardConfig
}
