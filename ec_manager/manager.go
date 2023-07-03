package ec_manager

import (
	"fmt"

	"github.com/BytemanD/ec-tools/common"
	"github.com/BytemanD/ec-tools/pkg/openstack/compute"
	"github.com/BytemanD/ec-tools/pkg/openstack/identity"
)

type ECManager struct {
	computeClient compute.ComputeClientV2
}

func (manager *ECManager) Init() error {
	if common.CONF.Auth.Url == "" {
		return fmt.Errorf("Auth.Url 未配置")
	}
	authClient, err := identity.GetV3Client(
		common.CONF.Auth.Url,
		common.CONF.Auth.User,
		common.CONF.Auth.Project,
		common.CONF.Auth.RegionName,
	)
	if err != nil {
		return fmt.Errorf("获取认证客户端失败, %s", err)
	}

	if err := authClient.TokenIssue(); err != nil {
		return fmt.Errorf("获取 Token 失败, %s", err)
	}
	manager.computeClient, err = compute.GetComputeClientV2(authClient)
	if err != nil {
		return fmt.Errorf("获取计算客户端失败, %s", err)
	}
	manager.computeClient.UpdateVersion()
	return nil
}
