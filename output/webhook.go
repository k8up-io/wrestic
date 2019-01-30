package output

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type webhookPusher struct {
	url string
}

func newWebhookHandler(url string) *webhookPusher {
	return &webhookPusher{
		url: url,
	}
}

func (w *webhookPusher) Push(hook JsonMarshaller) error {
	data := hook.ToJson()

	if len(data) <= 0 {
		return fmt.Errorf("webhook data is empty")
	}

	postBody := bytes.NewReader(hook.ToJson())

	resp, err := http.Post(w.url, "application/json", postBody)
	if err != nil || !strings.HasPrefix(resp.Status, "200") {
		httpCode := ""
		if resp == nil {
			httpCode = "http status unavailable"
		} else {
			httpCode = resp.Status
		}
		return fmt.Errorf("Could not send webhook: %v http status code: %v", err, httpCode)
	}
	return nil
}
