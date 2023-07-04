package httpclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/BytemanD/easygo/pkg/global/logging"
)

type Session struct {
}

func (session *Session) Request(req *http.Request, headers map[string]string) (Response, error) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	logging.Debug("Req: %s %s with headers: %v, body: %v", req.Method, req.URL, headers, req.Body)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	content, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	response := Response{
		Body:    content,
		Status:  resp.StatusCode,
		Headers: resp.Header}
	logging.Debug("Status: %s, Body: %s", resp.Status, content)
	return response, response.JudgeStatus()
}

func (session *Session) Get(url string, query map[string]string, headers map[string]string) (Response, error) {
	var params []string
	if query != nil {
		for k, v := range query {
			params = append(params, fmt.Sprintf("%s=%s", k, v))
		}
	}
	if len(params) != 0 {
		url += "?" + strings.Join(params, "&")
	}
	req, _ := http.NewRequest("GET", url, nil)
	return session.Request(req, headers)
}

func (session *Session) Post(url string, body []byte, headers map[string]string) (Response, error) {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	return session.Request(req, headers)
}

func (session *Session) Put(url string, body []byte, headers map[string]string) (Response, error) {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	return session.Request(req, headers)
}

func (session *Session) Delete(url string, headers map[string]string) (Response, error) {
	req, _ := http.NewRequest("DELETE", url, nil)
	return session.Request(req, headers)
}
