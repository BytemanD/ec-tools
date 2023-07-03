package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Status  int
	Body    []byte
	Headers http.Header
}

type HttpError struct {
	Code   int
	Reason string
}

func (err *HttpError) Error() string {
	return fmt.Sprintf("%d %s", err.Code, err.Reason)
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

func (resp *Response) BodyString() string {
	return string(resp.Body)
}

func (resp *Response) GetHeader(key string) string {
	return resp.Headers.Get(key)
}

func (resp *Response) BodyUnmarshal(object interface{}) error {
	return json.Unmarshal(resp.Body, object)
}
