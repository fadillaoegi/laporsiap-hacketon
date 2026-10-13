package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"laporsiap/internal/domain"
)

func (s *Store) ExportSnapshot(ctx context.Context, id, token string) (domain.Case, []domain.AttachmentData, error) {
	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return domain.Case{}, nil, err
	}
	defer tx.Rollback(ctx)
	var raw []byte
	err = tx.QueryRow(ctx, "SELECT document FROM cases WHERE id=$1 AND token_hash=$2 AND expires_at>now() AND access_expires_at>now()", id, domain.TokenHash(token)).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Case{}, nil, domain.ErrNotFound
	}
	if err != nil {
		return domain.Case{}, nil, err
	}
	var c domain.Case
	if err = json.Unmarshal(raw, &c); err != nil {
		return c, nil, err
	}
	rows, err := tx.Query(ctx, "SELECT id,filename,content_type,sha256,data FROM attachments WHERE case_id=$1 ORDER BY id", id)
	if err != nil {
		return c, nil, err
	}
	attachments := []domain.AttachmentData{}
	for rows.Next() {
		var a domain.AttachmentData
		if err = rows.Scan(&a.Meta.ID, &a.Meta.Filename, &a.Meta.ContentType, &a.Meta.SHA256, &a.Bytes); err != nil {
			rows.Close()
			return c, nil, err
		}
		a.Meta.Size = len(a.Bytes)
		attachments = append(attachments, a)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return c, nil, err
	}
	return c, attachments, tx.Commit(ctx)
}
