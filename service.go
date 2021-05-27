package kratos

// Service is an instance of a service in a discovery system.
type Service struct {
	ServiceID        string            `json:"id"`
	Servicename      string            `json:"name"`
	ServiceVersion   string            `json:"version"`
	ServiceMetadata  map[string]string `json:"metadata"`
	ServiceEndpoints []string          `json:"endpoints"`
}

// ID is service id
func (s *Service) ID() string {
	return s.ServiceID
}

// Name is service name
func (s *Service) Name() string {
	return s.Servicename
}

// Version is service Version
func (s *Service) Version() string {
	return s.ServiceVersion
}

// Metadata is service Metadata
func (s *Service) Metadata() map[string]string {
	return s.ServiceMetadata
}

// Endpoints is service Endpoints
func (s *Service) Endpoints() []string {
	return s.ServiceEndpoints
}
