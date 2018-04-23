package api

import (
	"fmt"
	"io"
	"mlog"
	"net/http"
	"time"
)

const defaultBufSize = 1024

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

	var buf []byte
	if resp.ContentLength > 0 {
		buf = make([]byte, resp.ContentLength)
	} else {
		buf = make([]byte, 0)
	}

	var (
		nerr  error
		readN = defaultBufSize
		read  int
	)

	if resp.ContentLength > 0 {
		_, nerr = io.ReadFull(resp.Body, buf)
	} else {
		nbuf := make([]byte, defaultBufSize)
		for readN == defaultBufSize && nerr == nil {
			readN, nerr = io.ReadFull(resp.Body, nbuf)
			read += readN
			buf = append(buf, nbuf...)
		}
		buf = buf[0:read]
	}

	if err != nil && err != io.ErrUnexpectedEOF {
		mlog.Error(fmt.Sprintf("[HTTP_GET]%s %s read:%s", err, resp.Status, buf))
		return nil, err
	}

	delay := (time.Now().Sub(start)) / time.Millisecond
	logInfo := fmt.Sprintf("[HTTP_GET]%s\r\n%dms %s", url, delay, string(buf))
	mlog.Info(logInfo)

	return buf, nil
}
