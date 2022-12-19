package data

type Country struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	Tlds     []string `json:"tld"`
	Capitals []string `json:"capital"`
	Flag     string   `json:"flag"`
	Maps     struct {
		GoogleMaps string `json:"googleMaps"`
	} `json:"maps"`
	Population int `json:"population"`
	Car        struct {
		Side string `json:"side"`
	} `json:"car"`
	Flags struct {
		Png string `json:"png"`
	} `json:"flags"`
}
