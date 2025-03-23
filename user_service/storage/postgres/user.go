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

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) storage.UserRepoI {
	return &UserRepo{
		db: db,
	}
}

// Create implements storage.UserRepoI.
func (s *UserRepo) Create(ctx context.Context, req *us.User) (*us.User, error) {
	id := uuid.NewString()

	_, err := s.db.Exec(ctx, `
		INSERT INTO users (
			id,
			user_type,
			user_role,
			full_name,
			user_name,
			email,
			password,
			gender,
			status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)`, id, req.UserType, req.UserRole, req.FullName, req.UserName, req.Email, req.Password, req.Gender, req.Status)

	if err != nil {
		log.Println("error while creating user in storage", err)
		return nil, err
	}

	user, err := s.GetSingle(ctx, &us.UserSingleRequest{Id: id})
	if err != nil {
		log.Println("error while getting user by id after creating", err)
		return nil, err
	}
	return user, nil
}

// GetByID implements storage.UserRepoI.
func (s *UserRepo) GetSingle(ctx context.Context, req *us.UserSingleRequest) (*us.User, error) {
	resp := &us.User{}

	var (
		created_at, updated_at time.Time
		query                  string
		args                   []interface{}
	)

	if req.Id != "" {
		query = `
			SELECT 
				id,
				user_type,
				user_role,
				full_name,
				user_name,
				email,
				password,
				gender,
				status,
				created_at,
				updated_at
			FROM users 
			WHERE id=$1`
		args = append(args, req.Id)
	} else if req.Email != "" {
		query = `
			SELECT 
				id,
				user_type,
				user_role,
				full_name,
				user_name,
				email,
				password,
				gender,
				status,
				created_at,
				updated_at
			FROM users 
			WHERE email=$1`
		args = append(args, req.Email)
	} else if req.Username != "" {
		query = `
			SELECT 
				id,
				user_type,
				user_role,
				full_name,
				user_name,
				email,
				password,
				gender,
				status,
				created_at,
				updated_at
			FROM users 
			WHERE user_name=$1`
		args = append(args, req.Username)
	} else {
		return nil, fmt.Errorf("either id, email, or username must be provided")
	}

	err := s.db.QueryRow(ctx, query, args...).Scan(&resp.Id, &resp.UserType, &resp.UserRole, &resp.FullName, &resp.UserName, &resp.Email, &resp.Password, &resp.Gender, &resp.Status, &created_at, &updated_at)

	if err != nil {
		if err == sql.ErrNoRows {
			return &us.User{}, fmt.Errorf("user not found")
		}
		log.Println("error while getting user", err)
		return &us.User{}, err
	}

	resp.CreatedAt = created_at.Format(time.RFC3339)
	resp.UpdatedAt = updated_at.Format(time.RFC3339)

	return resp, nil
}

// GetList implements storage.UserRepoI.
func (s *UserRepo) GetList(ctx context.Context, req *us.GetListUserRequest) (*us.GetListUserResponse, error) {
	resp := &us.GetListUserResponse{}
	var (
		filter                 = " WHERE TRUE"
		created_at, updated_at time.Time
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter += fmt.Sprintf(` AND (user_name ILIKE '%%%v%%' OR full_name ILIKE '%%%v%%' OR email ILIKE '%%%v%%')`, req.Search, req.Search, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)

	rows, err := s.db.Query(ctx, `
		SELECT
			id,
			user_type,
			user_role,
			full_name,
			user_name,
			email,
			gender,
			status,
			created_at,
			updated_at
		FROM users
	`+filter)

	if err != nil {
		log.Println("error while getting all users:", err)
		return nil, err
	}

	defer rows.Close()

	var count int64

	for rows.Next() {
		var user us.User
		count++
		err = rows.Scan(&user.Id, &user.UserType, &user.UserRole, &user.FullName, &user.UserName, &user.Email, &user.Gender, &user.Status, &created_at, &updated_at)

		if err != nil {
			log.Println("error while scanning users:", err)
			return nil, err
		}
		user.CreatedAt = created_at.Format(time.RFC3339)
		user.UpdatedAt = updated_at.Format(time.RFC3339)

		resp.Users = append(resp.Users, &user)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

// Update implements storage.UserRepoI.
func (s *UserRepo) Update(ctx context.Context, req *us.User) (*us.User, error) {

	_, err := s.db.Exec(ctx, `
        UPDATE users SET
		    full_name=$1,
		    user_name=$2,
			password=$3,
			gender=$4,
			status=$5,
            updated_at = NOW()
        WHERE id = $6`, req.FullName, req.UserName, req.Password, req.Gender, req.Status, req.Id)

	if err != nil {
		log.Println("error while updating user in storage", err)
		return nil, err
	}

	user, err := s.GetSingle(ctx, &us.UserSingleRequest{Id: req.Id})
	if err != nil {
		log.Println("error while getting updated user by id", err)
		return nil, err
	}

	return user, nil
}

// Delete implements storage.UserRepoI.
func (s *UserRepo) Delete(ctx context.Context, req *us.UserPrimaryKey) (*emptypb.Empty, error) {
	_, err := s.db.Exec(ctx, `
		DELETE FROM users
		WHERE id = $1
	`, req.Id)

	if err != nil {
		log.Println("error while deleting user:", err)
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
