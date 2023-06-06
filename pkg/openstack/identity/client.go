// Openstack 认证客户端

package identity

import (
	"fmt"
	"os"
	"strconv"
)

const DEFAULT_TOKEN_EXPIRE_SECOND = 3600

func getTokenExpireSecond() int {
	osSecond := os.Getenv("OS_TOKEN_EXPIRE_SECOND")
	if osSecond == "" {
		return DEFAULT_TOKEN_EXPIRE_SECOND
	} else {
		second, err := strconv.Atoi(osSecond)
		if err == nil {
			return second
		} else {
			return DEFAULT_TOKEN_EXPIRE_SECOND
		}

	}
}

// 初始化客户端前，先导入环境变量。
//
// 使用环境变量 'OS_TOKEN_EXPIRE_SECOND' 控制 Token 超时时间, 默认 3600s。
func GetV3ClientFromEnv() (V3AuthClient, error) {
	client := V3AuthClient{
		AuthUrl:           os.Getenv("OS_AUTH_URL"),
		Username:          os.Getenv("OS_USERNAME"),
		Password:          os.Getenv("OS_PASSWORD"),
		ProjectName:       os.Getenv("OS_PROJECT_NAME"),
		UserDomainName:    os.Getenv("OS_USER_DOMAIN_NAME"),
		ProjectDomainName: os.Getenv("OS_PROJECT_DOMAIN_NAME"),
		TokenExpireSecond: getTokenExpireSecond(),
	}
	if client.AuthUrl == "" {
		return client, fmt.Errorf("OS_AUTH_URL not found")
	}
	if client.AuthUrl == "" {
		return client, fmt.Errorf("OS_REGION_NAME not found")
	}
	return client, nil
}
