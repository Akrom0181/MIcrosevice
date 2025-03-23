package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	us "user_service/genproto/user_service"
	"user_service/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SessionRepo struct {
	db *pgxpool.Pool
}

func NewSessionRepo(db *pgxpool.Pool) storage.SessionRepoI {
	return &SessionRepo{
		db: db,
	}
}

// Create implements storage.SessionRepoI.
func (s *SessionRepo) Create(ctx context.Context, req *us.Session) (*us.Session, error) {
	id := uuid.NewString()

	expireDate := sql.NullTime{}
	expiresat, err := time.Parse(time.RFC3339, req.ExpiresAt)
	if err == nil {
		expireDate.Time = expiresat
		expireDate.Valid = true
	}

	_, err = s.db.Exec(ctx, `
        INSERT INTO session (
            id,
            user_id,
            ip_address,
            user_agent,
            is_active,
            expires_at,
            last_active_at,
            platform
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8
        )`, id, req.UserId, req.IpAddress, req.UserAgent, req.IsActive, expireDate, req.LastActiveAt, req.Platform)

	if err != nil {
		log.Println("error while creating session in storage", err)
		return nil, err
	}

	session, err := s.GetSingle(ctx, &us.SessionSingleRequest{Id: id})
	if err != nil {
		log.Println("error while getting session by id after creating", err)
		return nil, err
	}

	return session, nil
}

// GetByID implements storage.SessionRepoI.
func (s *SessionRepo) GetSingle(ctx context.Context, req *us.SessionSingleRequest) (*us.Session, error) {
	resp := &us.Session{}

	var (
		created_at, updated_at, last_active_at, expires_at time.Time
	)

	err := s.db.QueryRow(ctx, `
	        SELECT 
			id,
			user_id,
			ip_address,
			user_agent,
			is_active,
			expires_at,
			last_active_at,
			platform,
	        created_at,
	        updated_at
	        FROM session 
	    WHERE id=$1`, req.Id).Scan(
		&resp.Id,
		&resp.UserId,
		&resp.IpAddress,
		&resp.UserAgent,
		&resp.IsActive,
		&expires_at,
		&last_active_at,
		&resp.Platform,
		&created_at,
		&updated_at)

	if err != nil {
		log.Println("error while getting session by id", err)
		return nil, err
	}

	resp.CreatedAt = created_at.Format(time.RFC3339)
	resp.UpdatedAt = updated_at.Format(time.RFC3339)
	resp.LastActiveAt = last_active_at.Format(time.RFC3339)
	resp.ExpiresAt = expires_at.Format(time.RFC3339)

	return resp, nil
}

// GetList implements storage.SessionRepoI.
func (s *SessionRepo) GetList(ctx context.Context, req *us.GetListSessionRequest) (*us.GetListSessionResponse, error) {
	resp := &us.GetListSessionResponse{}
	var (
		filter                                             string
		created_at, updated_at, last_active_at, expires_at time.Time
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter = fmt.Sprintf(` WHERE user_id = '%s'`, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)

	query := `
        SELECT
            id,
			user_id,
			ip_address,
			user_agent,
			is_active,
			expires_at,
			last_active_at,
			platform,
	        created_at,
	        updated_at
        FROM session
    ` + filter

	rows, err := s.db.Query(ctx, query)

	if err != nil {
		log.Println("error while getting all sessions:", err)
		return nil, err
	}

	defer rows.Close()

	var count int64

	for rows.Next() {
		var user us.Session
		count++
		err = rows.Scan(
			&user.Id,
			&user.UserId,
			&user.IpAddress,
			&user.UserAgent,
			&user.IsActive,
			&expires_at,
			&last_active_at,
			&user.Platform,
			&created_at,
			&updated_at,
		)

		if err != nil {
			log.Println("error while scanning users:", err)
			return nil, err
		}
		user.CreatedAt = created_at.Format(time.RFC3339)
		user.UpdatedAt = updated_at.Format(time.RFC3339)
		user.LastActiveAt = last_active_at.Format(time.RFC3339)
		user.ExpiresAt = expires_at.Format(time.RFC3339)

		resp.Sessions = append(resp.Sessions, &user)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

// Update implements storage.SessionRepoI.
func (s *SessionRepo) Update(ctx context.Context, req *us.Session) (*us.Session, error) {

	_, err := s.db.Exec(ctx, `
        UPDATE session SET
		    ip_address=$1,
		    is_active=$2,
		    last_active_at=NOW(),
            updated_at = NOW()
        WHERE id = $3`, req.IpAddress, req.IsActive, req.Id)

	if err != nil {
		log.Println("error while updating session in storage", err)
		return nil, err
	}

	session, err := s.GetSingle(ctx, &us.SessionSingleRequest{Id: req.Id})
	if err != nil {
		log.Println("error while getting updated session by id", err)
		return nil, err
	}

	return session, nil
}

// Delete implements storage.SessionRepoI.
func (s *SessionRepo) Delete(ctx context.Context, req *us.SessionSingleRequest) (*emptypb.Empty, error) {
	_, err := s.db.Exec(ctx, `
		DELETE FROM session
		WHERE id = $1
	`, req.Id)

	if err != nil {
		log.Println("error while deleting session:", err)
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
