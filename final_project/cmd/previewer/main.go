package main

import (
	"flag"
	"fmt"

	"github.com/Lanworm/OTUS_GO/final_project/internal/config"
	"github.com/Lanworm/OTUS_GO/final_project/pkg/shortcuts"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./build/local/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := config.NewConfig(configFile)
	shortcuts.FatalIfErr(err)
	fmt.Println(config.Logger)
}
