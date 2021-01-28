package notification

import (
	"encoding/json"
	"github.com/lsgrep/sak/ulog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var client = &http.Client{}

type twilio struct {
	accountSid string
	authToken  string
	from       string
}

func NewTwilioApi(accountSid, authToken, from string) *twilio {
	return &twilio{accountSid: accountSid, authToken: authToken, from: from}
}

type twilioResponse struct {
	// meaning request has reached Twilio and has been queued
	Success bool `json:"success"`
	// in case, further checking delivery status etc are required
	ResponseData map[string]interface{} `json:"response_data"`
	ErrorMsg     string                 `json:"error_msg"`
}

func (t *twilio) SendSMS(to, body string) (*twilioResponse, error) {
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + t.accountSid + "/Messages.json"
	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", t.from)
	msgData.Set("Body", body)
	msgDataReader := *strings.NewReader(msgData.Encode())
	req, err := http.NewRequest("POST", urlStr, &msgDataReader)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(t.accountSid, t.authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	tr := &twilioResponse{}
	// if StatusCode is healthy, SMS is mostly probably has been queued
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		tr.Success = true
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		// parsing resp error is a minor issue
		if err != nil {
			ulog.Error(err)
			return tr, nil
		}
		tr.ResponseData = data
		return tr, nil
	} else {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ulog.Error(err)
			return tr, err
		}
		tr.ErrorMsg = string(bs)
		return tr, nil
	}
	return nil, nil
}
