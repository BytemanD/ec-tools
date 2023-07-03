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
	RegionName        string
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

func (authClient *V3AuthClient) TokenIssue() error {
	authBody := GetAuthReqBody(authClient.Username, authClient.Password, authClient.ProjectName)
	body, _ := json.Marshal(authBody)

	url := fmt.Sprintf("%s%s", authClient.AuthUrl, URL_AUTH_TOKEN)
	logging.Debug("Req: POST %s Body: %s", url, body)
	resp, err := http.Post(url, ContentType, bytes.NewBuffer(body))
	if err != nil {
		logging.Error("token issue failed, %s", err)
		return err
	}
	defer resp.Body.Close()
	logging.Debug("Status: %s", resp.Status)

	content, _ := ioutil.ReadAll(resp.Body)
	response := httpclient.Response{
		Status: resp.StatusCode, Headers: resp.Header, Body: content,
	}
	if err := response.JudgeStatus(); err != nil {
		logging.Error("token issue failed, %s", err)
		return err
	}

	var resToken RespToken
	json.Unmarshal(response.Body, &resToken)
	resToken.Token.tokenId = response.GetHeader("X-Subject-Token")
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

func (authClient *V3AuthClient) Request(method string, url string, body []byte, query map[string]string, headers map[string]string) (httpclient.Response, error) {
	var reqBody io.Reader = nil
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}
	req, _ := http.NewRequest(method, url, reqBody)
	tokenId := authClient.getTokenId()
	if tokenId == "" {
		return httpclient.Response{}, fmt.Errorf("token id is null")
	}
	req.Header.Set("X-Auth-Token", authClient.getTokenId())
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	logging.Debug("Req: %s %s with headers: %v, body: %v", method, url, headers, reqBody)
	resp, _ := http.DefaultClient.Do(req)
	content, _ := ioutil.ReadAll(resp.Body)
	logging.Debug("Status: %d, Body: %s", resp.StatusCode, content)
	defer resp.Body.Close()

	return httpclient.Response{Status: resp.StatusCode, Body: content}, nil
}

func (authClient *V3AuthClient) ServiceList() (httpclient.Response, error) {
	url := fmt.Sprintf("%s%s", authClient.AuthUrl, "/services")
	return authClient.Request("GET", url, nil, nil, nil)
}

func (authClient *V3AuthClient) UserList() (httpclient.Response, error) {
	url := fmt.Sprintf("%s%s", authClient.AuthUrl, "/users")
	return authClient.Request("GET", url, nil, nil, nil)
}

func (authClient *V3AuthClient) GetEndpointFromCatalog(serviceType string, endpointInterface string, region string) (string, error) {
	if len(authClient.token.Catalogs) == 0 {
		if err := authClient.TokenIssue(); err != nil {
			return "", err
		}
	}
	endpoints := authClient.token.GetEndpoints(OptionCatalog{
		Type:      serviceType,
		Interface: endpointInterface,
		Region:    region,
	})
	if (len(endpoints)) == 0 {
		return "", fmt.Errorf("endpoints not found")
	} else {
		return endpoints[0].Url, nil
	}
}
