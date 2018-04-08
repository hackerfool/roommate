package api

import (
	"fmt"
	"io"
	"mlog"
	"net/http"
	"time"
)

type httpHandle struct {
}

var httpHandler = &httpHandle{}

func (*httpHandle) Get(url string) ([]byte, error) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, buf)
	if err != nil {
		return nil, err
	}

	delay := (time.Now().Sub(start)) / time.Millisecond
	logInfo := fmt.Sprintf("[HTTP_GET]%s\r\n%dms %s", url, delay, string(buf))
	mlog.Info(logInfo)

	return buf, nil
}
