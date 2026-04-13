package model

// HealthStatus represents the payload returned by the health-check endpoint.
type HealthStatus struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Env     string `json:"env"`
}
