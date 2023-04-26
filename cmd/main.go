package main

import (
	"flag"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/sqlancer"
	"github.com/pingcap/tidb/util/signal"
	"github.com/sirupsen/logrus"
)

func flagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("gsqlancer", flag.ExitOnError)
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
	flagSet := flagSet()
	flagSet.Parse(os.Args[1:])
	configFile := flagSet.Lookup("config").Value.String()
	cfg, err := loadmeta(configFile)
	if err != nil {
		logrus.Error(err.Error())
		os.Exit(0)
	}
	svr := sqlancer.NewSQLancer(cfg)
	svr.Run()
	exited := make(chan struct{})
	signal.SetupSignalHandler(func(_ bool) {
		svr.Stop()
		close(exited)
	})
	<-exited
}
