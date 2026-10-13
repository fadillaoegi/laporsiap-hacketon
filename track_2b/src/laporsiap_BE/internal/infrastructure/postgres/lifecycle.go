package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"time"

	"github.com/jackc/pgx/v5"
	"laporsiap/internal/domain"
)

func (s *Store) checkCapacity(ctx context.Context, tx pgx.Tx, body []byte, id string, fileBytes int64, deleteID string) error {
	if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock(781234903)"); err != nil {
		return err
	}
	var count int
	var used, increase int64
	err := tx.QueryRow(ctx, `SELECT (SELECT count(*) FROM cases),
 (SELECT COALESCE(sum(octet_length(document::text)),0) FROM cases) +
 (SELECT COALESCE(sum(octet_length(document::text)),0) FROM case_history) +
 (SELECT COALESCE(sum(octet_length(data)),0) FROM attachments)`).Scan(&count, &used)
	if err != nil {
		return err
	}
	err = tx.QueryRow(ctx, `SELECT 2*octet_length($1::jsonb::text)::bigint -
 COALESCE((SELECT octet_length(document::text) FROM cases WHERE id=$2),0) + $3::bigint -
 COALESCE((SELECT octet_length(data) FROM attachments WHERE id=$4 AND case_id=$2),0)`, body, id, fileBytes, deleteID).Scan(&increase)
	if err != nil {
		return err
	}
	limits := s.PolicyLimits()
	if (id == "" && count >= limits.MaxCases) || used+increase > limits.StorageBytes {
		return domain.ErrCapacity
	}
	return nil
}

func (s *Store) DeleteCase(ctx context.Context, id, token string, revision int64) error {
	tag, err := s.Pool.Exec(ctx, `DELETE FROM cases WHERE id=$1 AND token_hash=$2 AND revision=$3 AND expires_at>now() AND access_expires_at>now()`, id, domain.TokenHash(token), revision)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (s *Store) RotateAccess(ctx context.Context, id, token string, revision int64) (domain.Case, string, error) {
	next := domain.RandomID()
	var raw []byte
	var expiry time.Time
	err := s.Pool.QueryRow(ctx, `UPDATE cases SET token_hash=$1,
 access_expires_at=LEAST(expires_at, now()+($2::bigint*interval '1 hour')),
 document=jsonb_set(document,'{access_expires_at}',to_jsonb(to_char(LEAST(expires_at,now()+($2::bigint*interval '1 hour')) AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS"Z"')))
 WHERE id=$3 AND token_hash=$4 AND revision=$5 AND expires_at>now() AND access_expires_at>now()
 RETURNING document, access_expires_at`, domain.TokenHash(next), s.PolicyLimits().AccessHours, id, domain.TokenHash(token), revision).Scan(&raw, &expiry)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Case{}, "", domain.ErrConflict
	}
	if err != nil {
		return domain.Case{}, "", err
	}
	var c domain.Case
	err = json.Unmarshal(raw, &c)
	c.AccessExpiresAt = expiry.UTC().Format(time.RFC3339)
	return c, next, err
}

// Automatic cleanup touches only expired cases. Downloads/backups are outside
// the database and cannot be recalled. Anonymous AI counts survive 48h.
func (s *Store) PurgeExpired(ctx context.Context) (int64, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	tag, err := tx.Exec(ctx, "DELETE FROM cases WHERE expires_at<=now()")
	if err != nil {
		return 0, err
	}
	if _, err = tx.Exec(ctx, "DELETE FROM ai_attempts WHERE created_at<now()-interval '48 hours'"); err != nil {
		return 0, err
	}
	return tag.RowsAffected(), tx.Commit(ctx)
}
func (s *Store) ExpiredCount(ctx context.Context) (int64, error) {
	var n int64
	err := s.Pool.QueryRow(ctx, "SELECT count(*) FROM cases WHERE expires_at<=now()").Scan(&n)
	return n, err
}
