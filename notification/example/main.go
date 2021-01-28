package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/lsgrep/sak/notification"
	"github.com/lsgrep/sak/ulog"
)

func main() {
	accountSid := "XXXXXX"
	authToken := "XXXXXXX"
	from := "+15092958538"
	twilio := notification.NewTwilioApi(accountSid, authToken, from)
	resp, err := twilio.SendSMS("+8618311191111", "Hello from DEx.top")
	if err != nil {
		ulog.Panic(err)
	}
	spew.Dump(resp)
}
