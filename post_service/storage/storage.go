package storage

import (
	"context"
	us "post_service/genproto/post_service"

	"google.golang.org/protobuf/types/known/emptypb"
)

type StorageI interface {
	CloseDB()
	PostAttachment() PostAttachmentRepoI
	Post() PostRepoI
}

type (

	// PostAttachmentRepoI -.
	PostAttachmentRepoI interface {
		Create(ctx context.Context, req *us.Attachment) (*us.Attachment, error)
		MultipleUpsert(ctx context.Context, req *us.AttachmentMultipleInsertRequest) (*us.AttachmentList, error)
		GetSingle(ctx context.Context, req *us.AttachmentSingleRequest) (*us.Attachment, error)
		GetList(ctx context.Context, req *us.GetListAttachmentRequest) (*us.AttachmentList, error)
		Delete(ctx context.Context, req *us.AttachmentSingleRequest) (*emptypb.Empty, error)
		GetDefaultTags(ctx context.Context, req *us.GetDefaultTagsRequest) (*us.GetDefaultTagsResponse, error)
	}

	// PostRepoI -.
	PostRepoI interface {
		Create(ctx context.Context, req *us.Post) (*us.Post, error)
		GetSingle(ctx context.Context, req *us.PostSingleRequest) (*us.Post, error)
		GetList(ctx context.Context, req *us.GetListPostRequest) (*us.PostList, error)
		Update(ctx context.Context, req *us.Post) (*us.Post, error)
		Delete(ctx context.Context, req *us.PostSingleRequest) (*emptypb.Empty, error)
	}
)
