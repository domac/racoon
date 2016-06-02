package main

import (
	"github.com/phillihq/racoon/app"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app.Main()
}
