package common

type Info struct {
	Version   string `json:"version"`
	GoVersion string `json:"go_version"`
	BuildId   string `json:"build_id"`
}

type Page struct {
	Page int `json:"page"`
	Size int `json:"size"`
}
