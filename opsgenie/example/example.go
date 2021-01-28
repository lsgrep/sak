package main

import (
	"fmt"
	"github.com/lsgrep/sak/opsgenie"
	"github.com/lsgrep/sak/ulog"
)

func main() {
	hb, err := opsgenie.NewHeartBeat("YOUR_API_KEY",
		"test_heartbeat",
		"test",
		5)
	if err != nil {
		ulog.Fatal(err)
	}
	hb.Start()

	// hb.ReportError(fmt.Errorf("this is test"), "demo of Alert API")
	fmt.Scanln()
}
