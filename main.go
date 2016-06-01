package main

import (
	"github.com/codegangsta/cli"
	"github.com/phillihq/racoon/util"
	"os"
	"runtime"
)

var configFileFlag = util.AddFlagString(cli.StringFlag{
	Name:   "config",
	EnvVar: "CONFIG",
	Value:  "config.json",
	Usage:  "the path of your config file",
})

//应用执行方法
func appAction(c *cli.Context) error {

	confileFilePath := c.String(configFileFlag.Name)

	println("-------", confileFilePath)

	//退出信号处理
	signalCH := util.InitSignal()
	util.HandleSignal(signalCH)
	return nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := cli.NewApp()
	app.Name = "racoon"
	app.Usage = "log collector, base on Go"
	app.Version = "0.0.1"
	app.Flags = util.GetAppFlags()
	app.Action = util.ActionWrapper(appAction)
	app.Run(os.Args)
}
