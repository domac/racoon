package exec

import (
	"bytes"
	"errors"
	"github.com/phillihq/racoon/config"
	"github.com/phillihq/racoon/config/gather"
	"os"
	osexec "os/exec"
	"strings"
	"time"
)

const ModuleName = "exec"

//输入配置
type InputConfig struct {
	config.InputConfig
	Command   string   `json:"command"`
	Args      []string `json:"args,omitempty"`
	Interval  int      `json:"interval,omitempty"`
	MsgTrim   string   `json:"message_trim,omitempty"`
	MsgPrefix string   `json:"message_prefix,omitempty"`

	hostname string `json:"-"`
}

func NewInputConfig() InputConfig {
	return InputConfig{
		InputConfig: config.InputConfig{
			StandardConfig: config.StandardConfig{
				Type: ModuleName,
			},
		},
		Interval: 60,
		MsgTrim:  "\t\r\n",
	}
}

func InitHandler(configitem *config.ConfigItem) (contextConfig config.InputContextConfig, err error) {
	conf := NewInputConfig()
	if err = config.ReflectConfig(configitem, &conf); err != nil {
		return
	}
	if conf.hostname, err = os.Hostname(); err != nil {
		return
	}
	contextConfig = &conf
	return
}

func (self *InputConfig) Start() {
	startChan := make(chan bool)
	ticker := time.NewTicker(time.Duration(self.Interval) * time.Second)

	go func() {
		startChan <- true
	}()

	for {
		select {
		case <-startChan:
			self.InvokeGet(self.Exec)
		case <-ticker.C:
			self.InvokeGet(self.Exec)
		}
	}
}

func (self *InputConfig) Exec(inchan config.InputCh) {
	errs := []error{}

	message, err := self.doExec()
	if err != nil {
		errs = append(errs, err)
	}

	extra := map[string]interface{}{
		"host": self.hostname,
	}

	gData := gather.GatherData{
		Message:   message,
		Timestamp: time.Now(),
		Extra:     extra,
	}

	if len(errs) > 0 {
		gData.AddTag("inputexec_failed")
	}
	inchan <- gData
	return
}

//命令执行
func (self *InputConfig) doExec() (data string, err error) {
	var (
		buffer bytes.Buffer
		raw    []byte
		cmd    *osexec.Cmd
	)

	cmd = osexec.Command(self.Command, self.Args...)
	cmd.Stderr = &buffer
	if raw, err = cmd.Output(); err != nil {
		return
	}
	data = string(raw)

	data = strings.Trim(data, self.MsgTrim)

	if buffer.Len() > 0 {
		err = errors.New(buffer.String())
	}
	return
}
