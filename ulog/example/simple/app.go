package main

/* Output:

180310-20:38:22.928 PST	[I]	app.go:24	hey
180310-20:38:22.928 PST	[E]	app.go:28	nicek1v1
180310-20:38:22.928 PST	[I]	app.go:32	match found!	{"txId": "0x1234abcd"}
180310-20:38:22.928 PST	[I]	app.go:37	old logger will work as expected!!!!

*/

import (
	"github.com/lsgrep/sak/ucfg"
	"github.com/lsgrep/sak/ulog"
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

func main() {
	logger := ulog.NewLogger()

	logger.Info("hey")
	logger.Debugw("nice", "k1", "v1", "k2", "v2")

	// error or above error level messages are also written to stderr
	logger.Error("nice", "k1", "v1")

	//contextual logging
	txLogger := logger.With("txId", "0x1234abcd")
	txLogger.Info("match found!")

	myValue := 100
	logger.Debugw("message", "my-key", myValue)

	logger.Info("old logger will work as expected!!!!")

	ulog.Info("Hey, this  is awesome.")
}
