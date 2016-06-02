package config

import (
	"errors"
	"fmt"
	"github.com/codegangsta/inject"
	"github.com/phillihq/racoon/config/gather"
	"github.com/phillihq/racoon/util"
)

type FilterContextConfig interface {
	ContextConfig
	Process(gather.GatherData) gather.GatherData
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
	filters, err := self.getFilters()
	if err != nil {
		return
	}
	go func() {
		for {
			select {
			case data := <-inchan:
				for _, filter := range filters {
					data = filter.Process(data)
				}
				outchan <- data
			}
		}
	}()
	return
}

func (self *Config) getFilters() (filters []FilterContextConfig, err error) {
	for _, configitem := range self.FilterItem {
		filterName := configitem["type"].(string)
		handler, ok := registedFilterHandlers[filterName]
		if !ok {
			err = errors.New(fmt.Sprintf("unknow filter config type:%s", filterName))
			return
		}

		inj := inject.New()
		inj.SetParent(self)
		inj.Map(&configitem)
		results, err := util.FuncInvoke(inj, handler)
		if err != nil {
			return []FilterContextConfig{}, err
		}

		for _, res := range results {
			if !res.CanInterface() {
				continue
			}

			if conf, ok := res.Interface().(FilterContextConfig); ok {
				conf.SetInjector(inj)
				filters = append(filters, conf)
			}
		}
	}
	return
}
