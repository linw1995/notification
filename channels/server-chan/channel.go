package server_chan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/linw1995/notification"
)

func New(endpoint, sendKey string) notification.Channel {
	return &Channel{
		Endpoint: endpoint,
		SendKey:  sendKey,
	}
}

type Channel struct {
	Endpoint string
	SendKey  string
}

type SendResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		PushId  string `json:"pushid"`
		ReadKey string `json:"readkey"`
		Error   string `json:"error"`
		Errno   int    `json:"errno"`
	}
}

func (ch *Channel) Send(ctx context.Context, opts ...notification.SendOption) (err error) {
	var (
		options = notification.GenerateOptions(opts...)
		cli     = http.DefaultClient
		sendUrl = fmt.Sprintf("%s/%s.send", ch.Endpoint, ch.SendKey)
		payload = make(url.Values)
	)

	if options.Title == "" {
		return fmt.Errorf("title is required")
	}
	payload.Add("title", options.Title)
	if options.Description != "" {
		payload.Add("desp", options.Description)
	}
	buf := bytes.NewReader([]byte(payload.Encode()))

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		sendUrl,
		buf,
	)
	if err != nil {
		err = fmt.Errorf("new request failed: %w", err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := cli.Do(req)
	if err != nil {
		err = fmt.Errorf("send request failed: %w", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result SendResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		err = fmt.Errorf("decode response failed: %w", err)
		return
	}
	if result.Code > 0 {
		err = fmt.Errorf("send failed %d: %s", resp.StatusCode, result.Message)
	}
	return
}
