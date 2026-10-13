package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"laporsiap/internal/domain"
	"time"
)

// Reserve before sending; rejected/failed provider requests count. Serializing
// reservations in PostgreSQL prevents bypass after restart or concurrent calls.
func (s *Store) ReserveAI(ctx context.Context, id, token string, revision int64) (string, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, "SELECT pg_advisory_xact_lock(781234904)"); err != nil {
		return "", err
	}
	var current int64
	err = tx.QueryRow(ctx, `SELECT revision FROM cases WHERE id=$1 AND token_hash=$2 AND expires_at>now() AND access_expires_at>now()`, id, domain.TokenHash(token)).Scan(&current)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", domain.ErrNotFound
	}
	if err != nil {
		return "", err
	}
	if current != revision {
		return "", domain.ErrConflict
	}
	var daily, global, concurrent int
	err = tx.QueryRow(ctx, `SELECT count(*) FILTER (WHERE case_id=$1 AND created_at>=now()-interval '24 hours'),
 count(*) FILTER (WHERE created_at>=now()-interval '24 hours'),
 count(*) FILTER (WHERE finished_at IS NULL AND lease_until>now()) FROM ai_attempts`, id).Scan(&daily, &global, &concurrent)
	if err != nil {
		return "", err
	}
	limits := s.PolicyLimits()
	if daily >= limits.AICaseDaily || global >= limits.AIGlobalDaily {
		return "", domain.ErrAIQuota
	}
	if concurrent >= limits.AIConcurrent {
		return "", domain.ErrAIBusy
	}
	attempt := domain.RandomID()
	_, err = tx.Exec(ctx, "INSERT INTO ai_attempts(id,case_id) VALUES($1,$2)", attempt, id)
	if err != nil {
		return "", err
	}
	return attempt, tx.Commit(ctx)
}
func (s *Store) FinishAI(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, _ = s.Pool.Exec(ctx, "UPDATE ai_attempts SET finished_at=now() WHERE id=$1", id)
}
