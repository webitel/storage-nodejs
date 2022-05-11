package microsoft

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	c := &client{
		key:    "0c50bb70b34e4bec9e68ccc587b79a18",
		region: "northeurope",
		http:   http.Client{},
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Minute)

	task, err := c.TranscriptJob("https://dev.webitel.com/api/storage/recordings/58600/stream?access_token=3gtm3xs4nbdi8ei59ihwcahzra", "uk-UA")
	if err != nil {
		t.Error(err.Error())
	}

	c.WaitFoSuccess(ctx, task)

	data, err := c.LoadTranscript(task)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(data))
}
