package identity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/BytemanD/ec-tools/pkg/httpclient"
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
func post(url string, body []byte) httpclient.Response {
	logging.Debug("Req: %s %s", url, body)
	resp, err := http.Post(url, ContentType, bytes.NewBuffer(body))
	if err != nil {
		return httpclient.Response{}
	}
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)
	xxxx := httpclient.Response{
		Status:  resp.StatusCode,
		Body:    content,
		Headers: resp.Header}
	logging.Debug("Resp: %s", xxxx.Body)
	return xxxx
}

func (authClient *V3AuthClient) TokenIssue() error {
	authBody := GetAuthReqBody(authClient.Username, authClient.Password, authClient.ProjectName)
	body, _ := json.Marshal(authBody)
	// TODO: use authClient.Request
	resp := post(fmt.Sprintf("%s%s", authClient.AuthUrl, URL_AUTH_TOKEN), body)
	if err := resp.JudgeStatus(); err != nil {
		return err
	}
	// var content, headers

	var resToken RespToken
	json.Unmarshal(resp.Body, &resToken)
	resToken.Token.tokenId = resp.GetHeader("X-Subject-Token")
	authClient.token = resToken.Token
	authClient.expiredAt = time.Now().Add(time.Second * time.Duration(authClient.TokenExpireSecond))
	return nil
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

func (authClient *V3AuthClient) Request(method string, url string, body []byte, query map[string]string, headers map[string]string) string {
	var reqBody io.Reader = nil
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}
	req, _ := http.NewRequest(method, url, reqBody)
	req.Header.Set("X-Auth-Token", authClient.getTokenId())
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	logging.Debug("Req: %s %s with %v", method, url, reqBody)
	resp, _ := http.DefaultClient.Do(req)
	content, _ := ioutil.ReadAll(resp.Body)
	logging.Debug("Body: %s", content)
	defer resp.Body.Close()

	return string(content)
}

func (authClient *V3AuthClient) ServiceList() string {
	url := fmt.Sprintf("%s%s", authClient.AuthUrl, "/services")
	return authClient.Request("GET", url, nil, nil, nil)
}

func (authClient *V3AuthClient) UserList() string {
	url := fmt.Sprintf("%s%s", authClient.AuthUrl, "/users")
	return authClient.Request("GET", url, nil, nil, nil)
}

func (authClient *V3AuthClient) GetEndpointFromCatalog(serviceType string, endpointInterface string, region string) string {
	if len(authClient.token.Catalogs) == 0 {
		authClient.TokenIssue()
	}
	endpoints := authClient.token.GetEndpoints(OptionCatalog{
		Type:      serviceType,
		Interface: endpointInterface,
		Region:    region,
	})
	if (len(endpoints)) == 0 {
		return ""
	} else {
		return endpoints[0].Url
	}
}
