package server_chan

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/linw1995/notification"
)

var (
	TestEndpoint string
	TestSendKey  string
)

func init() {
	TestEndpoint = os.Getenv("SERVER_CHAN_ENDPOINT")
	TestSendKey = os.Getenv("SERVER_CHAN_SEND_KEY")
}

func TestSend(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	ch := New(TestEndpoint, TestSendKey)
	err := ch.Send(
		ctx,
		notification.Title("通知测试"),
		notification.Description("这是一条测试通知"),
	)
	if err != nil {
		t.Fatal(err)
	}
}
