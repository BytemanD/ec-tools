package compute

import (
	"fmt"
	"os"

	"github.com/BytemanD/ec-tools/pkg/openstack/identity"
)

func GetComputeClientV2(authClient identity.V3AuthClient) (ComputeClientV2, error) {
	regionName := os.Getenv("OS_REGION_NAME")
	if regionName == "" {
		regionName = "RegionOne"
	}
	endpoint := authClient.GetEndpointFromCatalog(
		identity.TYPE_COMPUTE, identity.INTERFACE_PUBLIC, regionName)
	if endpoint == "" {
		return ComputeClientV2{}, fmt.Errorf("public endpoint for %s is not found", identity.TYPE_COMPUTE)
	}
	return ComputeClientV2{AuthClient: authClient, Endpoint: endpoint}, nil
}
