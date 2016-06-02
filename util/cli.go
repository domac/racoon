package util

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/codegangsta/inject"
	"log"
	"reflect"
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

//函数调用
func FuncInvoke(injector inject.Injector, f interface{}) (results []reflect.Value, err error) {
	if results, err = injector.Invoke(f); err != nil {
		return
	}
	for _, res := range results {
		if res.IsValid() {
			resI := res.Interface()
			switch resI.(type) {
			case error:
				err = resI.(error)
				break
			}
		}
	}
	return
}
