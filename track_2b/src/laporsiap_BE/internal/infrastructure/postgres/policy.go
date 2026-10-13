package postgres

import (
	"context"
	"laporsiap/internal/domain"
)

func (s *Store) PolicyLimits() domain.SafetyLimits {
	if s == nil || s.Limits.RetentionDays == 0 {
		return domain.DefaultSafetyLimits()
	}
	return s.Limits
}

func (s *Store) Ping(ctx context.Context) error { return s.Pool.Ping(ctx) }
