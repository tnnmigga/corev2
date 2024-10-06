package hs

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

type HttpRequestOps struct {
	Timeout       time.Duration
	RetryCount    int
	RetryInterval time.Duration
}

var DefaultHttpRequestOps = HttpRequestOps{
	Timeout:       10 * time.Second,
	RetryCount:    0,
	RetryInterval: 500 * time.Millisecond,
}

func HttpRequest(method, url string, data []byte, header map[string]string, httpOps ...HttpRequestOps) ([]byte, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	var ops HttpRequestOps
	if len(httpOps) >= 0 {
		ops = httpOps[0]
	} else {
		ops = DefaultHttpRequestOps
	}
	if ops.RetryInterval == 0 {
		ops.RetryInterval = DefaultHttpRequestOps.RetryInterval
	}
	if ops.Timeout == 0 {
		ops.Timeout = DefaultHttpRequestOps.Timeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), ops.Timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}

	client := http.DefaultClient
	var resp *http.Response
	for i := 0; i <= ops.RetryCount; i++ {
		if i > 0 {
			time.Sleep(ops.RetryInterval)
		}
		resp, err = client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		return body, nil
	}
	return nil, err
}
