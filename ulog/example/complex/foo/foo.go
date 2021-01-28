package foo

import (
	"github.com/lsgrep/sak/ulog"
	"sync"
)

func Run(wg *sync.WaitGroup) {
	logger := ulog.NewLoggerWithLevel("debug")

	var i = 0
	for i < 5 {
		logger.Debugw("some message!!!!", "idx", i)
		i++
	}
	wg.Done()
}
