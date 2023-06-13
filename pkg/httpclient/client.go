package httpclient

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/BytemanD/easygo/pkg/global/logging"
)

type Session struct {
}

func (session *Session) Request(method string, url string, body []byte, query map[string]string, headers map[string]string) (Response, error) {
	var reqBody io.Reader = nil
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}
	req, _ := http.NewRequest(method, url, reqBody)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	logging.Debug("Req: %s %s with %v", method, url, reqBody)
	resp, _ := http.DefaultClient.Do(req)
	content, _ := ioutil.ReadAll(resp.Body)
	// logging.Debug("Body: %s", content)
	defer resp.Body.Close()

	response := Response{
		Body:    content,
		Status:  resp.StatusCode,
		Headers: resp.Header}
	err := response.JudgeStatus()
	return response, err
}

func (session *Session) Get(url string, query map[string]string, headers map[string]string) (Response, error) {
	return session.Request("GET", url, nil, query, headers)
}

func (session *Session) Post(url string, body []byte, headers map[string]string) (Response, error) {
	return session.Request("POST", url, body, nil, headers)
}

func (session *Session) Put(url string, body []byte, headers map[string]string) (Response, error) {
	return session.Request("PUT", url, body, nil, headers)
}

func (session *Session) Delete(url string, query map[string]string, headers map[string]string) (Response, error) {
	return session.Request("DELETE", url, nil, nil, headers)
}
