package main

import (
	"goFoundation/webook/pkg/viperx"
	"log"
)

func main() {
	viperx.InitViperV1()
	initPrometheus()
	app := InitAPP()
	for _, v := range app.consumers {
		err := v.Start()
		if err != nil {
			panic(err)
		}
	}
	err := app.server.Serve()
	log.Println(err)
}
