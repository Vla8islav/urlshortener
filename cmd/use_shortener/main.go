package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"github.com/Vla8islav/urlshortener/internal/app/helpers"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

func main() {
	originalURL := helpers.GenerateRandomURL()

	resp := shortenURL(originalURL)

	shortenedURL := string(resp.Body())

	newCookie := resp.Header().Get("Set-Cookie")

	println(shortenedURL)
	println(newCookie)
}

func shortenURL(originalURL string) *resty.Response {
	errRedirectBlocked := errors.New("HTTP redirect blocked")
	redirPolicy := resty.RedirectPolicyFunc(func(_ *http.Request, _ []*http.Request) error {
		return errRedirectBlocked
	})

	httpc := resty.New().
		SetBaseURL("http://localhost:8889").
		SetRedirectPolicy(redirPolicy).
		SetProxy("http://localhost:8888")

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

	return resp
}

func deleteURL(originalURL string) *resty.Response {
	errRedirectBlocked := errors.New("HTTP redirect blocked")
	redirPolicy := resty.RedirectPolicyFunc(func(_ *http.Request, _ []*http.Request) error {
		return errRedirectBlocked
	})

	httpc := resty.New().
		SetBaseURL("http://localhost:8889").
		SetRedirectPolicy(redirPolicy).
		SetProxy("http://localhost:8888")

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

	return resp
}
