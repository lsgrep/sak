package opsgenie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lsgrep/sak/ulog"
	"io/ioutil"
	"net/http"
	"time"
)

const apiUrl = "https://api.opsgenie.com/v2"

type Heartbeat struct {
	ApiKey   string
	Name     string        // Name of the heartbeat
	Team     string        // Team name to assign
	Interval time.Duration // Interval in seconds to send heartbeat requests
	quit     chan int
}

func NewHeartBeat(apiKey string, name string, team string, interval time.Duration) (*Heartbeat, error) {
	if len(apiKey) == 0 {
		return nil, fmt.Errorf("missing OpsGenie apikey")
	}
	if len(team) == 0 && len(team) > 200 {
		team = "ops_team"
	}
	if interval < 1 {
		interval = 60
	}
	return &Heartbeat{
		ApiKey:   fmt.Sprintf("GenieKey %s", apiKey),
		Name:     name,
		Team:     team,
		Interval: interval * time.Second,
	}, nil
}

func sendHeartbeat(h *Heartbeat) error {
	url := fmt.Sprintf("%s/heartbeats/%s/ping", apiUrl, h.Name)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", h.ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		ulog.Errorf("Error sending heartbeat request %s", err.Error())
		return err
	}
	resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		body, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(body)
		ulog.Errorf("Error sending heartbeat %s", bodyString)
		return fmt.Errorf(bodyString)
	}
	ulog.Infof("Sent Heartbeat")
	return nil
}

type NewPing struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	IntervalUnit string `json:"intervalUnit"`
	Interval     int    `json:"interval"`
	Enabled      bool   `json:"enabled"`
	OwnerTeam    struct {
		Name string `json:"name"`
	} `json:"ownerTeam"`
}

func createHeartbeat(h *Heartbeat) error {
	newPing := NewPing{
		Name:         h.Name,
		Description:  "",
		IntervalUnit: "minutes",
		Interval:     5, // OpsGenie will pushing alert if missed heartbeat for 5 minutes
		Enabled:      true,
		OwnerTeam: struct {
			Name string `json:"name"`
		}{
			Name: h.Team,
		},
	}
	body, _ := json.Marshal(newPing)

	url := fmt.Sprintf("%s/heartbeats", apiUrl)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Add("Authorization", h.ApiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		ulog.Infof("Error creating heartbeat")
	} else {
		resp.Body.Close()
		if resp.StatusCode/100 != 2 && resp.StatusCode != 409 { // ignore 409 error if heartbeat already exists
			body, _ := ioutil.ReadAll(resp.Body)
			bodyString := string(body)
			ulog.Warningf("Error creating heartbeat[%d] %s", resp.StatusCode, bodyString)
			return fmt.Errorf(bodyString)
		}
		ulog.Infof("Created Heartbeat %s", h.Name)
	}
	return nil
}

type Alert struct {
	Message     string `json:"message"`
	Description string `json:"description"`
	Entity      string `json:"entity"`
	Priority    string `json:"priority"`
	Responders  []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"responders"`
}

// ReportError will instantly create a 'P1' alert in OpsGenie, use this only if it's emergent
func (h *Heartbeat) ReportError(reportErr error, desc string) {
	alert := Alert{
		Message:     reportErr.Error(),
		Description: desc,
		Entity:      h.Name,
		Priority:    "P1",
		Responders: []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		}{
			{
				Name: h.Team,
				Type: "team",
			},
		},
	}

	body, _ := json.Marshal(alert)
	url := fmt.Sprintf("%s/alerts", apiUrl)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Add("Authorization", h.ApiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		ulog.Errorf("Failed to sending alert %v. Got err %v", reportErr, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		body, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(body)
		ulog.Errorf("Failed to report to OpsGenie. OriginalErr: %v, Status: %d, Body: %s",
			reportErr.Error(), resp.StatusCode, bodyString)
	} else {
		ulog.Infof("Sent alert: %s", reportErr.Error())
	}
}

// Start sending heartbeat request.
func (h *Heartbeat) Start() error {
	err := createHeartbeat(h)
	if err != nil {
		return err
	}
	go sendHeartbeat(h)

	h.quit = make(chan int)

	go func(h *Heartbeat) {
		ticker := time.NewTicker(h.Interval)
		for {
			select {
			case <-ticker.C:
				sendHeartbeat(h)
			case <-h.quit:
				return
			}
		}
	}(h)

	return nil
}

func (h *Heartbeat) Stop() {
	h.quit <- 0
}
