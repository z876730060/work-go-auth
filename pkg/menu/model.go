package menu

type Menu struct {
	Key       string `json:"key"`
	Label     string `json:"label"`
	ParentKey string `json:"parentKey"`
}

type Route struct {
	Path      string `json:"path"`
	Component string `json:"component"`
	Other     bool   `json:"other"`
	MicroApp  string `json:"microApp"`
}
