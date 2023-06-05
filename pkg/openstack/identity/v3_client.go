package identity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/fjboy/magic-pocket/pkg/global/logging"
)

type V3AuthClient struct {
	AuthUrl           string
	Username          string
	Password          string
	ProjectName       string
	UserDomainName    string
	ProjectDomainName string
	TokenExpireSecond int
	token             Token
	expiredAt         time.Time
}

const (
	ContentType    string = "application/json"
	URL_AUTH_TOKEN string = "/auth/tokens"

	TYPE_COMPUTE  string = "compute"
	TYPE_IDENTITY string = "identity"

	INTERFACE_PUBLIC   string = "public"
	INTERFACE_ADMIN    string = "admin"
	INTERFACE_INTERVAL string = "internal"
)

// Return: body, headers
func post(url string, body []byte) ([]byte, http.Header) {
	logging.Debug("Req: %s", url)
	logging.Debug("Body: %s", body)
	resp, _ := http.Post(url, ContentType, bytes.NewBuffer(body))
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)
	logging.Debug("Resp: %s", content)
	return content, resp.Header
}

func (authClient *V3AuthClient) TokenIssue() {
	authBoy := GetAuthReqBody(authClient.Username, authClient.Password, authClient.ProjectName)
	body, _ := json.Marshal(authBoy)
	content, headers := post(fmt.Sprintf("%s%s", authClient.AuthUrl, URL_AUTH_TOKEN), body)

	var resToken RespToken
	json.Unmarshal(content, &resToken)
	resToken.Token.tokenId = headers.Get("X-Subject-Token")
	authClient.token = resToken.Token
	authClient.expiredAt = time.Now().Add(time.Second * time.Duration(authClient.TokenExpireSecond))
}
func (authClient *V3AuthClient) isTokenExpired() bool {
	if authClient.token.tokenId == "" {
		return true
	}
	if authClient.expiredAt.Before(time.Now()) {
		logging.Debug("token is exipred, expire second is %d", authClient.TokenExpireSecond)
		return true
	}
	return false
}

func (authClient *V3AuthClient) getTokenId() string {
	if authClient.isTokenExpired() {
		authClient.TokenIssue()
	}
	return authClient.token.tokenId
}

func (authClient *V3AuthClient) doRequest(method string, url string, body io.Reader) string {
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("X-Auth-Token", authClient.getTokenId())
	resp, _ := http.DefaultClient.Do(req)
	content, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return string(content)
}

func (authClient *V3AuthClient) ServiceList() string {
	url := fmt.Sprintf("%s%s", authClient.AuthUrl, "/services")
	return authClient.doRequest("GET", url, nil)
}

func (authClient *V3AuthClient) UserList() string {
	url := fmt.Sprintf("%s%s", authClient.AuthUrl, "/users")
	return authClient.doRequest("GET", url, nil)
}
