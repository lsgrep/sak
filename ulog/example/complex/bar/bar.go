package bar

import (
	"github.com/lsgrep/sak/ulog"
	"sync"
)

func Run(wg *sync.WaitGroup) {
	logger := ulog.NewLoggerWithLevel("error")
	var i = 0
	logger.Info("given current logger is `error` level, this message will be ignored")

	for i < 5 {
		logger.Error("caution, bar msg")
		i++
	}
	wg.Done()
}
