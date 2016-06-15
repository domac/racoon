package queue

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
	"os"
)

type logger interface {
	Output(int, string) error
}

type Consumer struct {
	client      *nsq.Consumer
	config      *nsq.Config
	nsqds       []string
	nsqlookupds []string
	concurrency int
	channel     string
	topic       string
	level       nsq.LogLevel
	log         logger
	err         error
}

//初始化消费端
func NewConsumer(topic, channel string) *Consumer {
	return &Consumer{
		log:         log.New(os.Stderr, "", log.LstdFlags),
		config:      nsq.NewConfig(),
		level:       nsq.LogLevelInfo,
		channel:     channel,
		topic:       topic,
		concurrency: 1,
	}
}

func (c *Consumer) SetLogger(log logger, level nsq.LogLevel) {
	c.level = level
	c.log = log
}

func (c *Consumer) SetMap(options map[string]interface{}) {
	for k, v := range options {
		c.Set(k, v)
	}
}

func (c *Consumer) Set(option string, value interface{}) {
	switch option {
	case "topic":
		c.topic = value.(string)
	case "channel":
		c.channel = value.(string)
	case "concurrency":
		c.concurrency = value.(int)
	case "nsqd":
		c.nsqds = []string{value.(string)}
	case "nsqlookupd":
		c.nsqlookupds = []string{value.(string)}
	case "nsqds":
		s, err := strings(value)
		if err != nil {
			c.err = fmt.Errorf("%q: %v", option, err)
			return
		}
		c.nsqds = s
	case "nsqlookupds":
		s, err := strings(value)
		if err != nil {
			c.err = fmt.Errorf("%q: %v", option, err)
			return
		}
		c.nsqlookupds = s
	default:
		err := c.config.Set(option, value)
		if err != nil {
			c.err = err
		}
	}
}

func (c *Consumer) Start(handler nsq.Handler) error {

	if c.err != nil {
		return c.err
	}

	client, err := nsq.NewConsumer(c.topic, c.channel, c.config)
	if err != nil {
		return err
	}
	c.client = client
	client.SetLogger(c.log, c.level)
	client.AddConcurrentHandlers(handler, c.concurrency)
	return c.connect()
}

//连接到nsqd
func (c *Consumer) connect() error {

	if len(c.nsqds) == 0 && len(c.nsqlookupds) == 0 {
		return fmt.Errorf(`at least one "nsqd" or "nsqlookupd" address must be configured`)
	}

	if len(c.nsqds) > 0 {
		err := c.client.ConnectToNSQDs(c.nsqds)
		if err != nil {
			return err
		}
	}
	if len(c.nsqlookupds) > 0 {
		err := c.client.ConnectToNSQLookupds(c.nsqlookupds)
		if err != nil {
			return err
		}
	}
	return nil
}

//stop and wait
func (c *Consumer) Stop() error {
	c.client.Stop()
	<-c.client.StopChan
	return nil
}

func strings(v interface{}) ([]string, error) {
	switch v.(type) {
	case []string:
		return v.([]string), nil
	case []interface{}:
		var ret []string
		for _, e := range v.([]interface{}) {
			s, ok := e.(string)
			if !ok {
				return nil, fmt.Errorf("string expected")
			}
			ret = append(ret, s)
		}
		return ret, nil
	default:
		return nil, fmt.Errorf("strings expected")
	}
}
