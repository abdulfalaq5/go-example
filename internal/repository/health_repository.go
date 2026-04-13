package repository

// HealthRepository is reserved for future data-layer health checks
// (e.g. pinging a database). Currently it is a no-op.
type HealthRepository interface {
	Ping() error
}

type healthRepository struct{}

// NewHealthRepository creates a concrete HealthRepository.
func NewHealthRepository() HealthRepository {
	return &healthRepository{}
}

// Ping is a placeholder that always succeeds. Replace with a real DB ping
// when a storage layer is added.
func (r *healthRepository) Ping() error {
	return nil
}
