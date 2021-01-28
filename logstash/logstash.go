package logstash

import (
	"fmt"
	"net"
	"time"
)

type Client struct {
	hostUrl string
	timeout time.Duration
}

// `timeout` is the time after which a write fails
func NewClient(hostUrl string, timeout time.Duration) *Client {
	return &Client{
		hostUrl: hostUrl,
		timeout: timeout,
	}
}

func (client *Client) Writeln(message string) error {
	conn, err := net.DialTimeout("tcp", client.hostUrl, client.timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	deadline := time.Now().Add(client.timeout)
	err = conn.SetDeadline(deadline)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(conn, message)
	return err
}
