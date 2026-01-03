package approval

import (
	"context"
	"database/sql"
	"errors"
	"encoding/json"
	"github.com/sarkozy543/approval-workflow-system/internal/types"
)
var ErrInvalidStatus = errors.New("request is not in PENDING status")

type Store struct {
	db *sql.DB
}

// NewStore: dışarıdan *sql.DB alır, bu sayede test etmesi kolay olur
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// GetAll: tüm approval request'leri getirir
func (s *Store) GetAll(ctx context.Context) ([]types.ApprovalRequest, error) {
	const query = `
		SELECT id, source_env, target_env, change_payload, status, requested_by, created_at, updated_at
		FROM approval_requests
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []types.ApprovalRequest

	for rows.Next() {
		var r types.ApprovalRequest
		if err := rows.Scan(
			&r.ID,
			&r.SourceEnv,
			&r.TargetEnv,
			&r.ChangePayload,
			&r.Status,
			&r.RequestedBy,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateInput: API'den gelecek body için basit input modeli
type CreateInput struct {
	SourceEnv     string `json:"source_env"`
	TargetEnv     string `json:"target_env"`
	ChangePayload json.RawMessage `json:"change_payload"`
	RequestedBy   string `json:"requested_by"`
}

// Create: yeni bir approval request ekler
func (s *Store) Create(ctx context.Context, in CreateInput) (*types.ApprovalRequest, error) {
	if in.SourceEnv == "" || in.TargetEnv == "" || len(in.ChangePayload) == 0 || in.RequestedBy == "" {
		return nil, errors.New("missing required fields")
	}

	// 1) Transaction başlat
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 2) approval_requests tablosuna insert et
	const insertRequest = `
		INSERT INTO approval_requests (
			source_env, target_env, change_payload, status, requested_by
		)
		VALUES ($1, $2, $3, 'PENDING', $4)
		RETURNING id, source_env, target_env, change_payload, status, requested_by, created_at, updated_at
	`

	var r types.ApprovalRequest

	err = tx.QueryRowContext(ctx, insertRequest,
		in.SourceEnv,
		in.TargetEnv,
		in.ChangePayload,
		in.RequestedBy,
	).Scan(
		&r.ID,
		&r.SourceEnv,
		&r.TargetEnv,
		&r.ChangePayload,
		&r.Status,
		&r.RequestedBy,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// 3) approval_logs tablosuna CREATED log'u ekle
	const insertLog = `
		INSERT INTO approval_logs (request_id, action, action_by, note)
		VALUES ($1, $2, $3, $4)
	`

	if _, err := tx.ExecContext(ctx, insertLog,
		r.ID,           // request_id
		"CREATED",      // action
		in.RequestedBy, // action_by
		"",             // note (şimdilik boş)
	); err != nil {
		return nil, err
	}

	// 4) Hepsi başarılı -> commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *Store) GetByID(ctx context.Context, id string) (*types.ApprovalRequest, error) {
	const query = `
		SELECT id, source_env, target_env, change_payload, status, requested_by, created_at, updated_at
		FROM approval_requests
		WHERE id = $1
	`

	var r types.ApprovalRequest

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&r.ID,
		&r.SourceEnv,
		&r.TargetEnv,
		&r.ChangePayload,
		&r.Status,
		&r.RequestedBy,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
func (s *Store) approveOrReject(
	ctx context.Context,
	id string,
	newStatus string,
	action string,
	actionBy string,
	note string,
) (*types.ApprovalRequest, error) {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	const updateQuery = `
		UPDATE approval_requests
		SET status = $1, updated_at = now()
		WHERE id = $2 AND status = 'PENDING'
		RETURNING id, source_env, target_env, change_payload, status, requested_by, created_at, updated_at
	`

	var r types.ApprovalRequest

	err = tx.QueryRowContext(ctx, updateQuery, newStatus, id).Scan(
		&r.ID,
		&r.SourceEnv,
		&r.TargetEnv,
		&r.ChangePayload,
		&r.Status,
		&r.RequestedBy,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidStatus
		}
		return nil, err
	}

	const insertLog = `
		INSERT INTO approval_logs (request_id, action, action_by, note)
		VALUES ($1, $2, $3, $4)
	`

	if _, err := tx.ExecContext(ctx, insertLog, r.ID, action, actionBy, note); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &r, nil
}


// Approve: request'i APPROVED yapar ve loglar
func (s *Store) Approve(ctx context.Context, id, actionBy, note string) (*types.ApprovalRequest, error) {
	return s.approveOrReject(ctx, id, "APPROVED", "APPROVED", actionBy, note)
}

// Reject: request'i REJECTED yapar ve loglar
func (s *Store) Reject(ctx context.Context, id, actionBy, note string) (*types.ApprovalRequest, error) {
	return s.approveOrReject(ctx, id, "REJECTED", "REJECTED", actionBy, note)
}
func (s *Store) GetLogsForRequest(ctx context.Context, requestID string) ([]types.ApprovalLog, error) {
	const query = `
		SELECT id, request_id, action, action_by, note, created_at
		FROM approval_logs
		WHERE request_id = $1
		ORDER BY created_at ASC
	`

	rows, err := s.db.QueryContext(ctx, query, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []types.ApprovalLog

	for rows.Next() {
		var l types.ApprovalLog
		if err := rows.Scan(
			&l.ID,
			&l.RequestID,
			&l.Action,
			&l.ActionBy,
			&l.Note,
			&l.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
