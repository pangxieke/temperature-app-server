package models

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	Timeout50ms  = time.Millisecond * 50
	Timeout200ms = time.Millisecond * 200
	Timeout350ms = time.Millisecond * 350
	Timeout500ms = time.Millisecond * 500
	Timeout15s   = time.Second * 15
)

func PostJson(url string, data []byte, timeout time.Duration,
	setter ...func(*http.Request)) ([]byte, http.Header, int, error) {
	return reqJson(http.MethodPost, url, data, timeout, setter...)
}

func reqJson(method string, url string, data []byte, timeout time.Duration,
	setters ...func(*http.Request)) ([]byte, http.Header, int, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err == nil {
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete {
			request.Header.Set("Content-Type", "application/json")
			request.Header.Add("Accept-Charset", "UTF-8")
		}
	}
	return reqInner(request, timeout, setters...)
}

func reqInner(req *http.Request, timeout time.Duration,
	setters ...func(*http.Request)) (body []byte, header http.Header, statusCode int, err error) {
	if req == nil {
		return
	}
	for _, setter := range setters {
		setter(req)
	}
	if timeout <= 0 {
		timeout = Timeout200ms
	}
	httpClient := &http.Client{
		Timeout: timeout,
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		fmt.Sprintf("post resp err=%s, body=%s", err, string(body))
		err = fmt.Errorf("post resp err=%s", err)
		return
	}

	header = res.Header
	statusCode = res.StatusCode
	return
}

func HttpPost(url string, body []byte, header map[string]string) (data []byte, err error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return
	}
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	client := http.Client{
		//Timeout: time.Duration(30000) * time.Millisecond,
	}
	resp, err := client.Do(req)
	fmt.Printf("HttpPost, url=%s, body=%s, header=%+v", url, body, header)
	if resp == nil {
		err = fmt.Errorf("post client, resp is nil")
		return
	}
	defer resp.Body.Close()

	if err != nil {
		err = fmt.Errorf("post %s, error=%s", url, err)
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("post error, resp statusCode=%d", resp.StatusCode)
		return
	}
	data, _ = ioutil.ReadAll(resp.Body)

	return
}