package servicesroute

type VersionResponse struct {
	Release   string `json:"release"`
	BuildDate string `json:"build_date"`
	GitHash   string `json:"git_hash"`
}
