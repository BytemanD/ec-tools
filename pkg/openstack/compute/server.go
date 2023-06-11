package compute

import "encoding/json"

func (cmpCli *ComputeClientV2) ServerList(query map[string]string) []Server {
	serversBody := ServersBody{}

	content := cmpCli.AuthClient.Request(
		"GET", cmpCli.getUrl("servers", "", query),
		nil, query, cmpCli.BaseHeaders)
	json.Unmarshal([]byte(content), &serversBody)
	return serversBody.Servers
}

func (computeClient *ComputeClientV2) ServerShow(id string) Server {
	content := computeClient.AuthClient.Request(
		"GET", computeClient.getUrl("servers", id, nil), nil, nil, computeClient.BaseHeaders)
	serverBody := ServerBody{}
	json.Unmarshal([]byte(content), &serverBody)
	return serverBody.Server
}

func (computeClient *ComputeClientV2) ServerDelete(id string) {
	computeClient.AuthClient.Request(
		"DELETE", computeClient.getUrl("servers", id, nil), nil, nil, computeClient.BaseHeaders)
}

func (computeClient *ComputeClientV2) ServerCreate() {
	computeClient.AuthClient.Request(
		"POST", computeClient.getUrl("servers", "", nil), nil, nil, computeClient.BaseHeaders)
}
