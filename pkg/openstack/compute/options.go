package compute

type ServerOpt struct {
	Flavor           string      `json:"flavorRef"`
	Image            string      `json:"imageRef"`
	Name             string      `json:"name"`
	Networks         interface{} `json:"networks"`
	AvailabilityZone string      `json:"availability_zone,omitempty"`
}
