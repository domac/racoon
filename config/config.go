package config

import (
	"io/ioutil"
)

//配置信息结构
type Config struct {
}

//从配置文件中载入配置信息
func LoadConfigFromFile(path string) (config Config, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	return LoadConfigFromData(data)
}

//从字节数据中载入配置信息
func LoadConfigFromData(data []byte) (config Config, err error) {
	return
}
