package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Status  int
	Body    io.Reader
	Headers http.Header
}

func (resp *Response) JudgeStatus() error {
	switch {
	case resp.Status <= 400:
		return nil
	case resp.Status == 400:
		return fmt.Errorf("BadRequest")
	case resp.Status == 404:
		return fmt.Errorf("NotFound")
	case resp.Status == 500:
		return fmt.Errorf("BadRequest")
	default:
		return fmt.Errorf("ErrorCode %d", resp.Status)
	}
}

func (resp *Response) BodyBytes() []byte {
	bytes, _ := ioutil.ReadAll(resp.Body)
	return bytes
}

func (resp *Response) BodyString() string {
	return string(resp.BodyBytes())
}

func (resp *Response) GetHeader(key string) string {
	return resp.Headers.Get(key)
}

func (resp *Response) BodyUnmarshal(object interface{}) error {
	return json.Unmarshal(resp.BodyBytes(), object)
}
