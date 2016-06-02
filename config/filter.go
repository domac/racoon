package config

type FilterContextConfig interface {
	ContextConfig
	Process(string) string
}

//过滤器配置信息
type FilterConfig struct {
	StandardConfig
}

type FilterHandler interface{}

var registedFilterHandlers = map[string]FilterHandler{}

//注册过滤器
func RegistFilterHandler(name string, handler FilterHandler) {
	registedFilterHandlers[name] = handler
}

func (self *Config) RunFilters() (err error) {
	return self.InvokeFunc(self.runFilters)
}

func (self *Config) runFilters(inchan InputCh, outchan OutputCh) (err error) {
	return
}
