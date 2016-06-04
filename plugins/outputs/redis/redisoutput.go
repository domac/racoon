package redis

import (
	"errors"
	"fmt"
	"github.com/fzzy/radix/redis"
	"github.com/phillihq/racoon/config"
	"github.com/phillihq/racoon/config/gather"
	"time"
)

const ModuleName = "redis"

//Redis输出配置
type RedisOutputConfig struct {
	config.OutputConfig
	key               string   `json:"key"`
	Host              []string `json:"host"`
	DataType          string   `json:"data_type,omitempty"`
	Timeout           int      `json:"timeout,omitempty"`
	ReconnectInterval int      `json:"reconnect_interval,omitempty"`

	sendChan chan gather.GatherData
	client   *redis.Client
	clients  []*redis.Client
}

func NewRedisOutputConfig() RedisOutputConfig {
	return RedisOutputConfig{
		OutputConfig: config.OutputConfig{
			StandardConfig: config.StandardConfig{
				Type: ModuleName,
			},
		},
		key:               "racoon",
		DataType:          "list",
		Timeout:           5,
		ReconnectInterval: 1,
		sendChan:          make(chan gather.GatherData),
	}
}

func InitHandler(configitem *config.ConfigItem) (contextConfig config.OutputContextConfig, err error) {
	conf := NewRedisOutputConfig()
	if err = config.ReflectConfig(configitem, &conf); err != nil {
		return
	}
	go conf.loop()
	if err = conf.initRedisClient(); err != nil {
		return
	}
	contextConfig = &conf
	return
}

//初始化redis客户端
func (self *RedisOutputConfig) initRedisClient() (err error) {

	var client *redis.Client
	self.closeClients()

	for _, addr := range self.Host {
		if client, err = redis.DialTimeout("tcp", addr, time.Duration(self.Timeout)*time.Second); err == nil {
			self.clients = append(self.clients, client)
		} else {
			fmt.Errorf("Redis connection failed: %q\n%s", addr, err)
		}
	}

	if len(self.clients) > 0 {
		self.client = self.clients[0]
		err = nil
	} else {
		self.client = nil
		err = errors.New("no valid redis server connection")
	}
	return
}

//关闭客户端
func (self *RedisOutputConfig) closeClients() (err error) {

	var client *redis.Client
	for _, client = range self.clients {
		client.Close()
	}
	self.clients = self.clients[:0]
	return
}

func (self *RedisOutputConfig) loop() (err error) {

	for {
		data := <-self.sendChan
		self.sendData(data)
	}
	return
}

func (self *RedisOutputConfig) Send(data gather.GatherData) (err error) {
	self.sendChan <- data
	return
}

func (self *RedisOutputConfig) sendData(data gather.GatherData) (err error) {

	var (
		client *redis.Client
		raw    []byte
		key    string
	)

	if raw, err = data.MarshalJSON(); err != nil {
		fmt.Errorf("event Marshal failed: %v", data)
		return
	}

	key = self.key

	if self.client != nil {
		if err = self.redisSend(self.client, key, raw); err == nil {
			return
		}
	}

	for {
		if err = self.initRedisClient(); err != nil {
			return
		}

		for _, client = range self.clients {
			if err = self.redisSend(client, key, raw); err == nil {
				self.client = client
				return
			}
		}

		time.Sleep(time.Duration(self.ReconnectInterval) * time.Second)
	}

	return
}

func (self *RedisOutputConfig) redisSend(client *redis.Client, key string, raw []byte) (err error) {
	var res *redis.Reply
	switch self.DataType {
	case "list":
		res = client.Cmd("rpush", key, raw)
		err = res.Err
	case "channel":
		res = client.Cmd("publish", key, raw)
		err = res.Err
	default:
		err = errors.New("unknown DataType: " + self.DataType)
	}
	return
}
