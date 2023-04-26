package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/sirupsen/logrus"
)

func redpFlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("gsqlancer", flag.ExitOnError)
	flagSet.Bool("version", false, "print version string")
	flagSet.String("config", "", "path to config file")
	return flagSet
}

func loadmeta(configFile string) (meta *config.Config, err error) {
	if configFile != "" {
		_, err = toml.DecodeFile(configFile, &meta)
		if err != nil {
			return
		}
	}
	return
}

func main() {
	flagSet := redpFlagSet()
	flagSet.Parse(os.Args[1:])
	if flagSet.Lookup("version").Value.(flag.Getter).Get().(bool) || len(os.Args) == 1 {
		fmt.Println(version.String())
		os.Exit(0)
	}
	configFile := flagSet.Lookup("config").Value.String()
	config, err := loadmeta(configFile)
	if err != nil {
		logrus.Error(err.Error())
		os.Exit(0)
	}
}
