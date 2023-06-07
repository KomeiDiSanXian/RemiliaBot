// Package web 插件网络请求
package web

import (
	"errors"
	"net/http"
	"time"
)

// MakeRequest 用于获取网页信息，header自定义
//
// Client 超时时间为 30s
// 如果请求回复码不为200，会至多重复请求 3 次
func MakeRequest(url string, headers map[string]string) (*http.Response, error) {
	cli := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	var resp *http.Response
	for i := 0; i < 3; i++ {
		resp, err = cli.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == http.StatusOK {
			break
		}
	}

	if resp == nil {
		return nil, errors.New("response is nil")
	}

	return resp, nil
}
