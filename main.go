package main

import (
	"fmt"
	"github.com/RayofLightz/ingestion/extract"
	"github.com/RayofLightz/ingestion/streamproc"
)

func main() {
	config_set, err := extract.JsonConf()
	if err != nil {
		fmt.Println(err)
	} else {
		streamproc.StartProcessor(config_set)
	}
}
