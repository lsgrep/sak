package main

import (
	"github.com/lsgrep/sak/ucfg"
	"github.com/lsgrep/sak/ulog"
	"github.com/lsgrep/sak/ulog/example/complex/bar"
	"github.com/lsgrep/sak/ulog/example/complex/foo"
	"sync"

	"github.com/spf13/viper"
)

func init() {
	// config file setup
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// add current working directory and load config file
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	ucfg.Bootstrap()
}

var wg sync.WaitGroup

func main() {
	logger := ulog.NewLogger()

	logger.Info("hey")
	logger.Debug("nice", "k1", "v1", "k2", "v2")

	wg.Add(1)
	go foo.Run(&wg)

	wg.Add(1)
	go bar.Run(&wg)

	wg.Wait()
}
