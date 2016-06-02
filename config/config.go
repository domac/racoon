package config

import (
	"encoding/json"
	"github.com/codegangsta/inject"
	"github.com/phillihq/racoon/util"
	"io/ioutil"
	"reflect"
)

var bufferChannelSize = 100 //初始化缓冲通道大小

type InputCh chan string
type OutputCh chan string

//上下文配置接口
type ContextConfig interface {
	SetInjector(injector inject.Injector)
	GetType() string
	InvokeGet(f interface{}) (refvs []reflect.Value, err error)
}

type StandardConfig struct {
	inject.Injector `json:"-"`
	Type            string `json:"type"`
}

func (self *StandardConfig) SetInjector(injector inject.Injector) {
	self.Injector = injector
}

//获取类型
func (self *StandardConfig) GetType() string {
	return self.Type
}

//函数调用,并获取结果
func (self *StandardConfig) InvokeGet(f interface{}) (refvs []reflect.Value, err error) {
	inj := self.Injector
	return util.FuncInvoke(inj, f)

}

type ConfigItem map[string]interface{}

//配置信息结构
type Config struct {
	inject.Injector `json:"-"`
	InputItem       []ConfigItem `json:"input,omitempty"`
	OutputItem      []ConfigItem `json:"output,omitempty"`
	FilterItem      []ConfigItem `json:"filter,omitempty"`
}

//只调用,不获取结果
func (self *Config) InvokeFunc(f interface{}) error {
	_, err := util.FuncInvoke(self.Injector, f)
	return err
}

//从配置文件中载入配置信息
func LoadConfigFromFile(path string) (config Config, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	return LoadConfigFromData(data)
}

//从字符串中载入配置信息
func LoadConfigFromString(text string) (config Config, err error) {
	return LoadConfigFromData([]byte(text))
}

//从字节数据中载入配置信息
func LoadConfigFromData(data []byte) (config Config, err error) {

	if err = json.Unmarshal(data, &config); err != nil {
		return
	}
	config.Injector = inject.New()
	inputChannel := make(InputCh, bufferChannelSize)
	outputChannel := make(OutputCh, bufferChannelSize)

	//注入相关结构到配置上下文
	config.Map(inputChannel)
	config.Map(outputChannel)

	return
}

func ReflectConfig(configItem *ConfigItem, conf interface{}) (err error) {
	data, err := json.Marshal(configItem)
	if err != nil {
		return
	}

	if err = json.Unmarshal(data, conf); err != nil {
		return
	}

	return
}
