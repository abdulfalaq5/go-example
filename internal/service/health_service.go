package service

import (
	"github.com/falaqmsi/go-example/internal/model"
	"github.com/falaqmsi/go-example/internal/repository"
)

// HealthService defines the business logic for the health-check feature.
type HealthService interface {
	Check(env string) (*model.HealthStatus, error)
}

type healthService struct {
	repo repository.HealthRepository
}

// NewHealthService wires the service to its repository dependency.
func NewHealthService(repo repository.HealthRepository) HealthService {
	return &healthService{repo: repo}
}

// Check verifies service health and returns a structured status report.
func (s *healthService) Check(env string) (*model.HealthStatus, error) {
	if err := s.repo.Ping(); err != nil {
		return nil, err
	}

	return &model.HealthStatus{
		Status:  "ok",
		Version: "1.0.0",
		Env:     env,
	}, nil
}
