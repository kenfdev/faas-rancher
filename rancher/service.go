package rancher

// LaunchConfig refers to the rancher service's launch config
type LaunchConfig struct {
	Environment   map[string]string `json:"environment"`
	Labels        map[string]string `json:"labels"`
	RestartPolicy map[string]string `json:"restartPolicy"`
	ImageUUID     string            `json:"imageUuid"`
}

// ServiceResponse is the response structure for service requests
type ServiceResponse struct {
	Data []Service `json:"data"`
}

// Service refers to rancher's Service
type Service struct {
	ID            string        `json:"id"`
	StackID       string        `json:"stackId"`
	StartOnCreate bool          `json:"startOnCreate"`
	Name          string        `json:"name"`
	Scale         uint64        `json:"scale"`
	LaunchConfig  *LaunchConfig `json:"launchConfig"`
	State         string        `json:"state"`
}

// IsActive tells whether the service is active or not
func (s *Service) IsActive() bool {
	return s.State == "active"
}
