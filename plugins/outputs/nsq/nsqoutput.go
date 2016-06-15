package nsq

import (
	"github.com/nsqio/go-nsq"
	"github.com/phillihq/racoon/config"
	"github.com/phillihq/racoon/config/gather"
)

const (
	ModuleName = "nsq"
)

//NSQ输出配置
type NsqOutOutConfig struct {
	config.OutputConfig
	Nsqd  string `json:"nsqd"`
	Topic string `json:"topic"`

	producer *nsq.Producer
}

func NewNsqOutputConfig() NsqOutOutConfig {
	return NsqOutOutConfig{
		OutputConfig: config.OutputConfig{
			StandardConfig: config.StandardConfig{
				Type: ModuleName,
			},
		},
	}
}

//初始化处理
func InitHandler(configitem *config.ConfigItem) (contextConfig config.OutputContextConfig, err error) {
	conf := NewNsqOutputConfig()
	if err = config.ReflectConfig(configitem, &conf); err != nil {
		return
	}
	if err = conf.initNsqProducer(); err != nil {
		return
	}
	contextConfig = &conf
	return
}

//初始化NSQ的生产者
func (self *NsqOutOutConfig) initNsqProducer() (err error) {
	p, err := nsq.NewProducer(self.Nsqd, nsq.NewConfig())
	if err != nil {
		return err
	}
	self.producer = p
	err = nil
	return
}

//发送采集的数据到NSQ
func (self *NsqOutOutConfig) Send(data gather.GatherData) (err error) {
	msg := data.Message
	self.producer.Publish(self.Topic, []byte(msg))
	return
}
