package compute

type Flavor struct {
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	OriginalName string            `json:"original_name"`
	Ram          int               `json:"ram"`
	Vcpus        int               `json:"vcpus"`
	ExtraSpecs   map[string]string `json:"extra_specs"`
}

type Fault struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Details string `json:"details"`
}

type Server struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	VmState    string `json:"OS-EXT-STS:vm_state"`
	PowerState string `json:"OS-EXT-STS:power_state"`
	Host       string `json:"OS-EXT-SRV-ATTR:host"`
	AZ         string `json:"OS-EXT-AZ:availability_zone"`
	Flavor     Flavor `json:"flavor"`
	Fault      Fault  `json:"fault"`
}

type ServerBody struct {
	Server Server `json:"server"`
}

type ServersBody struct {
	Servers []Server `json:"servers"`
}

type ServeCreaterBody struct {
	Server ServerOpt `json:"server"`
}
