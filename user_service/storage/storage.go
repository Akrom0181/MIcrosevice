package storage

import (
	"context"
	us "user_service/genproto/user_service"

	"google.golang.org/protobuf/types/known/emptypb"
)

type StorageI interface {
	CloseDB()
	User() UserRepoI
	Session() SessionRepoI
}

type (
	UserRepoI interface {
		Create(ctx context.Context, req *us.User) (*us.User, error)
		GetSingle(ctx context.Context, req *us.UserSingleRequest) (*us.User, error)
		GetList(ctx context.Context, req *us.GetListUserRequest) (*us.GetListUserResponse, error)
		Update(ctx context.Context, req *us.User) (*us.User, error)
		Delete(ctx context.Context, req *us.UserPrimaryKey) (*emptypb.Empty, error)
	}

	SessionRepoI interface {
		Create(ctx context.Context, req *us.Session) (*us.Session, error)
		GetSingle(ctx context.Context, req *us.SessionSingleRequest) (*us.Session, error)
		GetList(ctx context.Context, req *us.GetListSessionRequest) (*us.GetListSessionResponse, error)
		Update(ctx context.Context, req *us.Session) (*us.Session, error)
		Delete(ctx context.Context, req *us.SessionSingleRequest) (*emptypb.Empty, error)
	}
)
