package postgres

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"laporsiap/internal/domain"
	"strings"
	"time"
)

//go:embed migrations/001_initial.sql
var migration string

//go:embed migrations/002_safety.sql
var safetyMigration string

type Store struct {
	Pool   *pgxpool.Pool
	Limits domain.SafetyLimits
}

func OpenStore(ctx context.Context, url string) (*Store, error) {
	if strings.TrimSpace(url) == "" {
		return nil, errors.New("DATABASE_URL wajib diisi secara eksplisit; koneksi database default tidak digunakan.")
	}
	p, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	if err = p.Ping(ctx); err != nil {
		p.Close()
		return nil, err
	}
	tx, err := p.Begin(ctx)
	if err != nil {
		p.Close()
		return nil, err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, "SELECT pg_advisory_xact_lock(781234901)"); err == nil {
		_, err = tx.Exec(ctx, migration+"\n"+safetyMigration)
	}
	if err == nil {
		err = tx.Commit(ctx)
	}
	if err != nil {
		p.Close()
		return nil, err
	}
	return &Store{Pool: p}, nil
}
func (s *Store) Create(ctx context.Context, c domain.Case, token string) (domain.Case, error) {
	body, _ := json.Marshal(c)
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return c, err
	}
	defer tx.Rollback(ctx)
	if err = s.checkCapacity(ctx, tx, body, "", 0, ""); err != nil {
		return c, err
	}
	if _, err = tx.Exec(ctx, "INSERT INTO cases(id,token_hash,document,expires_at,access_expires_at) VALUES($1,$2,$3,$4,$5)", c.ID, domain.TokenHash(token), body, c.ExpiresAt, c.AccessExpiresAt); err != nil {
		return c, err
	}
	if _, err = tx.Exec(ctx, "INSERT INTO case_history(case_id,revision,action,document) VALUES($1,1,'create',$2)", c.ID, body); err != nil {
		return c, err
	}
	return c, tx.Commit(ctx)
}
func (s *Store) Get(ctx context.Context, id, token string) (domain.Case, error) {
	var b []byte
	var expires, accessExpires time.Time
	err := s.Pool.QueryRow(ctx, "SELECT document, expires_at, access_expires_at FROM cases WHERE id=$1 AND token_hash=$2 AND expires_at>now() AND access_expires_at>now()", id, domain.TokenHash(token)).Scan(&b, &expires, &accessExpires)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Case{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Case{}, err
	}
	var c domain.Case
	err = json.Unmarshal(b, &c)
	c.ExpiresAt = expires.UTC().Format(time.RFC3339)
	c.AccessExpiresAt = accessExpires.UTC().Format(time.RFC3339)
	return c, err
}
func (s *Store) Save(ctx context.Context, c domain.Case, token, action string, expectedRevision int64, attach *domain.AttachmentData, deleteID string) (domain.Case, error) {
	old := expectedRevision
	if old >= s.PolicyLimits().MaxRevisions {
		return c, domain.ErrCapacity
	}
	body, err := json.Marshal(c)
	if err != nil {
		return c, err
	}
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return c, err
	}
	defer tx.Rollback(ctx)
	var fileBytes int64
	if attach != nil {
		fileBytes = int64(len(attach.Bytes))
	}
	if err = s.checkCapacity(ctx, tx, body, c.ID, fileBytes, deleteID); err != nil {
		return c, err
	}
	tag, err := tx.Exec(ctx, "UPDATE cases SET revision=$1,document=$2,updated_at=now() WHERE id=$3 AND token_hash=$4 AND revision=$5 AND expires_at>now() AND access_expires_at>now()", c.Revision, body, c.ID, domain.TokenHash(token), old)
	if err != nil {
		return c, err
	}
	if tag.RowsAffected() != 1 {
		return c, domain.ErrConflict
	}
	if attach != nil {
		a := attach.Meta
		if _, err = tx.Exec(ctx, "INSERT INTO attachments(id,case_id,filename,content_type,sha256,data) VALUES($1,$2,$3,$4,$5,$6)", a.ID, c.ID, a.Filename, a.ContentType, a.SHA256, attach.Bytes); err != nil {
			return c, err
		}
	}
	if deleteID != "" {
		if _, err = tx.Exec(ctx, "DELETE FROM attachments WHERE id=$1 AND case_id=$2", deleteID, c.ID); err != nil {
			return c, err
		}
	}
	if _, err = tx.Exec(ctx, "INSERT INTO case_history(case_id,revision,action,document) VALUES($1,$2,$3,$4)", c.ID, c.Revision, action, body); err != nil {
		return c, err
	}
	return c, tx.Commit(ctx)
}

func (s *Store) Attachment(ctx context.Context, caseID, id, token string) (domain.AttachmentData, error) {
	var a domain.AttachmentData
	err := s.Pool.QueryRow(ctx, `SELECT a.id,a.filename,a.content_type,a.sha256,a.data FROM attachments a JOIN cases c ON c.id=a.case_id WHERE a.case_id=$1 AND a.id=$2 AND c.token_hash=$3 AND c.expires_at>now() AND c.access_expires_at>now()`, caseID, id, domain.TokenHash(token)).Scan(&a.Meta.ID, &a.Meta.Filename, &a.Meta.ContentType, &a.Meta.SHA256, &a.Bytes)
	if errors.Is(err, pgx.ErrNoRows) {
		return a, domain.ErrNotFound
	}
	a.Meta.Size = len(a.Bytes)
	return a, err
}

func (s *Store) History(ctx context.Context, id, token string) ([]domain.HistoryEntry, error) {
	if _, err := s.Get(ctx, id, token); err != nil {
		return nil, err
	}
	rows, err := s.Pool.Query(ctx, "SELECT h.revision,h.action,h.created_at::text,h.document FROM case_history h JOIN cases c ON c.id=h.case_id WHERE c.id=$1 AND c.token_hash=$2 AND c.expires_at>now() AND c.access_expires_at>now() ORDER BY h.revision DESC LIMIT 50", id, domain.TokenHash(token))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.HistoryEntry{}
	for rows.Next() {
		var h domain.HistoryEntry
		if err = rows.Scan(&h.Revision, &h.Action, &h.CreatedAt, &h.Document); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, rows.Err()
}
