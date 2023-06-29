package compute

type ServerOpt struct {
	Flavor   string `json:"flavorRef"`
	Image    string `json:"imageRef"`
	Name     string `json:"name"`
	Networks string `json:"networks"`
}
