package config

import (
	"fmt"
	"github.com/codegangsta/inject"
	"github.com/phillihq/racoon/config/gather"
	"github.com/phillihq/racoon/util"
)

type OutputContextConfig interface {
	ContextConfig
	Send(gather.GatherData) (err error)
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
	outputs, err := self.getOutputs()
	if err != nil {
		return
	}
	go func() {
		for {
			select {
			case data := <-outchan:
				for _, output := range outputs {
					go func(o OutputContextConfig, data gather.GatherData) {
						if err = o.Send(data); err != nil {
							fmt.Errorf("output failed: %v\n", err)
						}
					}(output, data)
				}
			}
		}
	}()
	return
}

func (self *Config) getOutputs() (outputs []OutputContextConfig, err error) {

	for _, configitem := range self.OutputItem {
		outputName := configitem["type"].(string)
		handler, ok := registedOutputHandlers[outputName]
		if !ok {
			return
		}

		inj := inject.New()
		inj.SetParent(self)
		inj.Map(&configitem)
		results, err := util.FuncInvoke(inj, handler)
		if err != nil {
			return []OutputContextConfig{}, err
		}

		for _, res := range results {
			if !res.CanInterface() {
				continue
			}
			if conf, ok := res.Interface().(OutputContextConfig); ok {
				conf.SetInjector(inj)
				outputs = append(outputs, conf)
			}
		}
	}

	return
}
