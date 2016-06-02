package config

import (
	"errors"
	"fmt"
	"github.com/codegangsta/inject"
	"github.com/phillihq/racoon/util"
)

//输入配置接口
type InputContextConfig interface {
	ContextConfig
	Start()
}

type InputConfig struct {
	StandardConfig
}

type InputHandler interface{}

var registedInputHandlers = map[string]InputHandler{}

//注册
func RegistInputHandler(name string, handler InputHandler) {
	registedInputHandlers[name] = handler
}

func (self *Config) RunInputs() error {
	return self.InvokeFunc(self.runInputs)
}

func (self *Config) runInputs(inChan InputCh) error {

	inputs, err := self.getInputs(inChan)
	if err != nil {
		return err
	}

	for _, input := range inputs {
		go input.Start()
	}

	return nil
}

func (self *Config) getInputs(inChan InputCh) (inputs []InputContextConfig, err error) {
	for _, confItem := range self.InputItem {
		handler, ok := registedInputHandlers[confItem["type"].(string)]
		if !ok {
			err = errors.New(fmt.Sprintf("unknow input config type:%s", confItem["type"].(string)))
			return
		}

		inj := inject.New()
		inj.SetParent(self)
		inj.Map(&confItem)
		inj.Map(inChan)
		results, err := util.FuncInvoke(inj, handler)
		if err != nil {
			return []InputContextConfig{}, err
		}

		for _, res := range results {
			if !res.CanInterface() {
				continue
			}
			if conf, ok := res.Interface().(InputContextConfig); ok {
				conf.SetInjector(inj)
				inputs = append(inputs, conf)
			}
		}
	}
	return
}
