package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

func main() {
	shortenURLs := make(map[string]string)

	errRedirectBlocked := errors.New("HTTP redirect blocked")
	redirPolicy := resty.RedirectPolicyFunc(func(_ *http.Request, _ []*http.Request) error {
		return errRedirectBlocked
	})

	httpc := resty.New().
		SetBaseURL("http://localhost:8889").
		SetRedirectPolicy(redirPolicy).
		SetProxy("http://localhost:8888")

	originalURL := "http://ayaginkdkzmu.net/keu3mjdqmlun/jucsjdybso6s0"

	// сжимаем данные с помощью gzip
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, _ = zw.Write([]byte(originalURL))
	_ = zw.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// выполняем запрос с выставлением необходимых заголовков
	req := httpc.R().
		SetContext(ctx).
		SetBody(buf.Bytes()).
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Encoding", "gzip")
	resp, err := req.Post("/")
	if err != nil {
		println(err)
	}

	shortenURL := string(resp.Body())
	shortenURLs[originalURL] = shortenURL
}
