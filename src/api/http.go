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
	mlog.Debug("recv ", url, " data:", resp.ContentLength, " ", resp.Status)
	buf := make([]byte, 102400)
	_, err = io.ReadFull(resp.Body, buf)
	if err != nil {
		mlog.Error(fmt.Sprintf("[HTTP_GET]%s %s read:%s", err, resp.Status, buf))
		return nil, err
	}

	delay := (time.Now().Sub(start)) / time.Millisecond
	logInfo := fmt.Sprintf("[HTTP_GET]%s\r\n%dms %s", url, delay, string(buf))
	mlog.Info(logInfo)

	return buf, nil
}
