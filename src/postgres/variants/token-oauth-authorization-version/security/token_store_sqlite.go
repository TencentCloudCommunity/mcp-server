package security

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type SQLiteTokenStore struct {
	db *sql.DB
}

func NewSQLiteTokenStoreFromEnv() (*SQLiteTokenStore, error) {
	storeKind := strings.TrimSpace(os.Getenv(TokenStoreEnv))
	if storeKind != "" && !strings.EqualFold(storeKind, "sqlite") {
		return nil, fmt.Errorf("unsupported %s=%q", TokenStoreEnv, storeKind)
	}
	return NewSQLiteTokenStore(TokenStorePathFromEnv())
}

func NewSQLiteTokenStore(path string) (*SQLiteTokenStore, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		path = defaultTokenStorePath
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create token store dir: %w", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite token store: %w", err)
	}
	store := &SQLiteTokenStore{db: db}
	if err := store.initSchema(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *SQLiteTokenStore) CreateToken(ctx context.Context, record *TokenRecord) error {
	if record == nil {
		return fmt.Errorf("record is nil")
	}
	scopesJSON, err := json.Marshal(record.Scopes)
	if err != nil {
		return err
	}
	regionsJSON, err := json.Marshal(record.AllowedRegions)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `
INSERT INTO issued_tokens (
    id, token_hash, token_prefix, subject_type, subject_id, tenant_id,
    display_name, scopes_json, allowed_regions_json, status,
    expires_at, created_at, updated_at, last_used_at,
    issued_by, revoked_by, revoked_at, revoke_reason, description
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`,
		record.ID,
		record.TokenHash,
		record.TokenPrefix,
		record.SubjectType,
		record.SubjectID,
		record.TenantID,
		record.DisplayName,
		string(scopesJSON),
		string(regionsJSON),
		string(record.Status),
		record.ExpiresAt.Format(time.RFC3339Nano),
		record.CreatedAt.Format(time.RFC3339Nano),
		record.UpdatedAt.Format(time.RFC3339Nano),
		nullableTime(record.LastUsedAt),
		record.IssuedBy,
		nullableString(record.RevokedBy),
		nullableTime(record.RevokedAt),
		nullableString(record.RevokeReason),
		record.Description,
	)
	if err != nil {
		return fmt.Errorf("insert token: %w", err)
	}
	return nil
}

func (s *SQLiteTokenStore) GetTokenByID(ctx context.Context, id string) (*TokenRecord, error) {
	return s.queryOne(ctx, `SELECT id, token_hash, token_prefix, subject_type, subject_id, tenant_id, display_name, scopes_json, allowed_regions_json, status, expires_at, created_at, updated_at, last_used_at, issued_by, revoked_by, revoked_at, revoke_reason, description FROM issued_tokens WHERE id = ?`, id)
}

func (s *SQLiteTokenStore) GetTokenByHash(ctx context.Context, tokenHash string) (*TokenRecord, error) {
	return s.queryOne(ctx, `SELECT id, token_hash, token_prefix, subject_type, subject_id, tenant_id, display_name, scopes_json, allowed_regions_json, status, expires_at, created_at, updated_at, last_used_at, issued_by, revoked_by, revoked_at, revoke_reason, description FROM issued_tokens WHERE token_hash = ?`, tokenHash)
}

func (s *SQLiteTokenStore) ListTokens(ctx context.Context, filter TokenFilter) ([]TokenRecord, error) {
	query := `SELECT id, token_hash, token_prefix, subject_type, subject_id, tenant_id, display_name, scopes_json, allowed_regions_json, status, expires_at, created_at, updated_at, last_used_at, issued_by, revoked_by, revoked_at, revoke_reason, description FROM issued_tokens WHERE 1=1`
	args := make([]any, 0, 4)
	if tenantID := strings.TrimSpace(filter.TenantID); tenantID != "" {
		query += ` AND tenant_id = ?`
		args = append(args, tenantID)
	}
	if subjectID := strings.TrimSpace(filter.SubjectID); subjectID != "" {
		query += ` AND subject_id = ?`
		args = append(args, subjectID)
	}
	if status := strings.TrimSpace(string(filter.Status)); status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	query += ` ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list tokens: %w", err)
	}
	defer rows.Close()
	result := make([]TokenRecord, 0)
	for rows.Next() {
		record, err := scanTokenRecord(rows)
		if err != nil {
			return nil, err
		}
		if scope := strings.TrimSpace(filter.Scope); scope != "" {
			if !record.ToPrincipal().HasScope(scope) {
				continue
			}
		}
		result = append(result, *record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SQLiteTokenStore) RevokeToken(ctx context.Context, id string, input RevokeTokenInput) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	res, err := s.db.ExecContext(ctx, `
UPDATE issued_tokens
SET status = ?, updated_at = ?, revoked_by = ?, revoked_at = ?, revoke_reason = ?
WHERE id = ?
`, string(TokenStatusRevoked), now, nullableString(strings.TrimSpace(input.RevokedBy)), now, nullableString(strings.TrimSpace(input.Reason)), id)
	if err != nil {
		return fmt.Errorf("revoke token: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrTokenNotFound
	}
	return nil
}

func (s *SQLiteTokenStore) TouchToken(ctx context.Context, id string, usedAt time.Time) error {
	_, err := s.db.ExecContext(ctx, `UPDATE issued_tokens SET last_used_at = ?, updated_at = ? WHERE id = ?`, usedAt.Format(time.RFC3339Nano), usedAt.Format(time.RFC3339Nano), id)
	return err
}

func (s *SQLiteTokenStore) PutCredentialBinding(ctx context.Context, binding *CredentialBinding) error {
	if binding == nil {
		return fmt.Errorf("credential binding is nil")
	}
	createdAt := binding.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	updatedAt := binding.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	_, err := s.db.ExecContext(ctx, `
INSERT INTO issued_token_credentials (
    token_id, subject_id, provider, credential_kind, encrypted_body, expires_at, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(token_id) DO UPDATE SET
    subject_id = excluded.subject_id,
    provider = excluded.provider,
    credential_kind = excluded.credential_kind,
    encrypted_body = excluded.encrypted_body,
    expires_at = excluded.expires_at,
    updated_at = excluded.updated_at
`,
		binding.TokenID,
		binding.SubjectID,
		defaultString(strings.TrimSpace(binding.Provider), "tencentcloud"),
		strings.TrimSpace(binding.CredentialKind),
		binding.EncryptedBody,
		nullableTime(binding.ExpiresAt),
		createdAt.Format(time.RFC3339Nano),
		updatedAt.Format(time.RFC3339Nano),
	)
	if err != nil {
		return fmt.Errorf("upsert credential binding: %w", err)
	}
	return nil
}

func (s *SQLiteTokenStore) GetCredentialBinding(ctx context.Context, tokenID string) (*CredentialBinding, error) {
	row := s.db.QueryRowContext(ctx, `SELECT token_id, subject_id, provider, credential_kind, encrypted_body, expires_at, created_at, updated_at FROM issued_token_credentials WHERE token_id = ?`, tokenID)
	binding, err := scanCredentialBinding(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCredentialBindingNotFound
		}
		return nil, err
	}
	return binding, nil
}

func (s *SQLiteTokenStore) DeleteCredentialBinding(ctx context.Context, tokenID string) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM issued_token_credentials WHERE token_id = ?`, tokenID)
	if err != nil {
		return fmt.Errorf("delete credential binding: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrCredentialBindingNotFound
	}
	return nil
}

func (s *SQLiteTokenStore) initSchema(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS issued_tokens (
    id TEXT PRIMARY KEY,
    token_hash TEXT NOT NULL UNIQUE,
    token_prefix TEXT NOT NULL,
    subject_type TEXT NOT NULL,
    subject_id TEXT NOT NULL,
    tenant_id TEXT NOT NULL,
    display_name TEXT NOT NULL,
    scopes_json TEXT NOT NULL,
    allowed_regions_json TEXT NOT NULL,
    status TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    last_used_at TEXT NULL,
    issued_by TEXT NOT NULL,
    revoked_by TEXT NULL,
    revoked_at TEXT NULL,
    revoke_reason TEXT NULL,
    description TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_issued_tokens_subject ON issued_tokens(subject_id);
CREATE INDEX IF NOT EXISTS idx_issued_tokens_tenant ON issued_tokens(tenant_id);
CREATE INDEX IF NOT EXISTS idx_issued_tokens_status ON issued_tokens(status);
CREATE TABLE IF NOT EXISTS issued_token_credentials (
    token_id TEXT PRIMARY KEY,
    subject_id TEXT NOT NULL,
    provider TEXT NOT NULL,
    credential_kind TEXT NOT NULL,
    encrypted_body TEXT NOT NULL,
    expires_at TEXT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_issued_token_credentials_subject ON issued_token_credentials(subject_id);
`)
	if err != nil {
		return fmt.Errorf("init token schema: %w", err)
	}
	return nil
}

func (s *SQLiteTokenStore) queryOne(ctx context.Context, query string, args ...any) (*TokenRecord, error) {
	row := s.db.QueryRowContext(ctx, query, args...)
	record, err := scanTokenRecord(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}
	return record, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanTokenRecord(scanner rowScanner) (*TokenRecord, error) {
	var (
		record             TokenRecord
		scopesJSON         string
		allowedRegionsJSON string
		status             string
		lastUsedAt         sql.NullString
		revokedBy          sql.NullString
		revokedAt          sql.NullString
		revokeReason       sql.NullString
		expiresAt          string
		createdAt          string
		updatedAt          string
	)
	if err := scanner.Scan(
		&record.ID,
		&record.TokenHash,
		&record.TokenPrefix,
		&record.SubjectType,
		&record.SubjectID,
		&record.TenantID,
		&record.DisplayName,
		&scopesJSON,
		&allowedRegionsJSON,
		&status,
		&expiresAt,
		&createdAt,
		&updatedAt,
		&lastUsedAt,
		&record.IssuedBy,
		&revokedBy,
		&revokedAt,
		&revokeReason,
		&record.Description,
	); err != nil {
		return nil, err
	}
	record.Status = TokenStatus(status)
	if err := json.Unmarshal([]byte(scopesJSON), &record.Scopes); err != nil {
		return nil, fmt.Errorf("decode scopes: %w", err)
	}
	if err := json.Unmarshal([]byte(allowedRegionsJSON), &record.AllowedRegions); err != nil {
		return nil, fmt.Errorf("decode allowed regions: %w", err)
	}
	var err error
	record.ExpiresAt, err = time.Parse(time.RFC3339Nano, expiresAt)
	if err != nil {
		return nil, err
	}
	record.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return nil, err
	}
	record.UpdatedAt, err = time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return nil, err
	}
	if lastUsedAt.Valid {
		parsed, err := time.Parse(time.RFC3339Nano, lastUsedAt.String)
		if err != nil {
			return nil, err
		}
		record.LastUsedAt = &parsed
	}
	if revokedBy.Valid {
		record.RevokedBy = revokedBy.String
	}
	if revokedAt.Valid {
		parsed, err := time.Parse(time.RFC3339Nano, revokedAt.String)
		if err != nil {
			return nil, err
		}
		record.RevokedAt = &parsed
	}
	if revokeReason.Valid {
		record.RevokeReason = revokeReason.String
	}
	return &record, nil
}

func scanCredentialBinding(scanner rowScanner) (*CredentialBinding, error) {
	var (
		binding   CredentialBinding
		expiresAt sql.NullString
		createdAt string
		updatedAt string
	)
	if err := scanner.Scan(
		&binding.TokenID,
		&binding.SubjectID,
		&binding.Provider,
		&binding.CredentialKind,
		&binding.EncryptedBody,
		&expiresAt,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}
	var err error
	binding.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return nil, err
	}
	binding.UpdatedAt, err = time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return nil, err
	}
	if expiresAt.Valid {
		parsed, err := time.Parse(time.RFC3339Nano, expiresAt.String)
		if err != nil {
			return nil, err
		}
		binding.ExpiresAt = &parsed
	}
	return &binding, nil
}

func nullableTime(value *time.Time) any {
	if value == nil || value.IsZero() {
		return nil
	}
	return value.Format(time.RFC3339Nano)
}

func nullableString(value string) any {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return value
}
