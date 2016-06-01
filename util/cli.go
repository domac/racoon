package util

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
)

var appFlags = map[string]cli.Flag{}

//添加字符串参数
func AddFlagString(sf cli.StringFlag) cli.StringFlag {
	if _, ok := appFlags[sf.Name]; ok {
		panic(fmt.Sprintf("flag %s denined", sf.Name))
	} else {
		appFlags[sf.Name] = sf
	}
	return sf
}

func GetAppFlags() (afs []cli.Flag) {
	for _, f := range appFlags {
		afs = append(afs, f)
	}
	return
}

func ActionWrapper(action func(context *cli.Context) error) func(context *cli.Context) {
	return func(context *cli.Context) {
		if err := action(context); err != nil {
			log.Println(err.Error())
		}
	}
}
