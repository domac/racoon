package app

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/phillihq/racoon/config"
	"github.com/phillihq/racoon/util"
	"os"
)

var configFileFlag = util.AddFlagString(cli.StringFlag{
	Name:   "config",
	EnvVar: "CONFIG",
	Value:  "config.json",
	Usage:  "the path of your config file",
})

//应用执行方法
func appAction(c *cli.Context) (err error) {

	confileFilePath := c.String(configFileFlag.Name)

	//读取配置信息
	conf, err := config.LoadConfigFromFile(confileFilePath)
	if err != nil {
		return errors.New(fmt.Sprintf("load config file failed, %v", err))
	}

	if err = conf.RunInputs(); err != nil {
		return
	}

	if err = conf.RunFilters(); err != nil {
		return
	}

	if err = conf.RunOutputs(); err != nil {
		return
	}

	//退出信号处理
	signalCH := util.InitSignal()
	util.HandleSignal(signalCH)
	return
}

func Main() {
	app := cli.NewApp()
	app.Name = "racoon"
	app.Usage = "log collector, base on Go"
	app.Version = "0.0.1"
	app.Flags = util.GetAppFlags()
	app.Action = util.ActionWrapper(appAction)
	app.Run(os.Args)
}
