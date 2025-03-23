package postgres

import (
	"context"
	"fmt"
	"log"
	us "post_service/genproto/post_service"
	"post_service/storage"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PostAttachmentRepo struct {
	db *pgxpool.Pool
}

func NewAttachmentRepo(db *pgxpool.Pool) storage.PostAttachmentRepoI {
	return &PostAttachmentRepo{
		db: db,
	}
}

func (s *PostAttachmentRepo) MultipleUpsert(ctx context.Context, req *us.AttachmentMultipleInsertRequest) (*us.AttachmentList, error) {

	sql := `INSERT INTO post_attachment (id, post_id, filepath, content_type) VALUES`
	args := []interface{}{}
	argIndex := 1

	for _, attachment := range req.Attachments {
		if attachment.Id == "" {
			attachment.Id = uuid.NewString()
		}
		sql += fmt.Sprintf(" ($%d, $%d, $%d, $%d),", argIndex, argIndex+1, argIndex+2, argIndex+3)
		args = append(args, attachment.Id, req.PostId, attachment.Filepath, attachment.ContentType)
		argIndex += 4
	}

	sql = strings.TrimSuffix(sql, ",")

	_, err := s.db.Exec(ctx, sql, args...)
	if err != nil {
		log.Println("error while inserting multiple post_attachments", err)
		return nil, err
	}

	attachments, err := s.GetList(ctx, &us.GetListAttachmentRequest{
		Search: req.PostId,
		Limit:  100,
		Page:   1})
	if err != nil {
		log.Println("error while retrieving inserted post_attachments", err)
		return nil, err
	}

	return attachments, nil
}

func (s *PostAttachmentRepo) Create(ctx context.Context, req *us.Attachment) (*us.Attachment, error) {
	id := uuid.NewString()

	_, err := s.db.Exec(ctx, `
		INSERT INTO post_attachment (
			id,
			post_id,
			filepath,
			content_type
		) VALUES (
			$1, $2, $3, $4
		)`, id, req.PostId, req.Filepath, req.ContentType)

	if err != nil {
		log.Println("error while creating post_attachment in storage", err)
		return nil, err
	}

	attachment, err := s.GetSingle(ctx, &us.AttachmentSingleRequest{Id: id})
	if err != nil {
		log.Println("error while getting post_attachment by id after creating", err)
		return nil, err
	}
	return attachment, nil
}

func (s *PostAttachmentRepo) GetSingle(ctx context.Context, req *us.AttachmentSingleRequest) (*us.Attachment, error) {
	resp := &us.Attachment{}

	var (
		created_at, updated_at time.Time
	)

	err := s.db.QueryRow(ctx, `
	        SELECT 
			id,
			post_id,
			filepath,
			content_type,
	        created_at,
	        updated_at
	        FROM post_attachment 
	    WHERE id=$1`, req.Id).Scan(&resp.Id, &resp.PostId, &resp.Filepath, &resp.ContentType, &created_at, &updated_at)

	if err != nil {
		log.Println("error while getting attachment by id", err)
		return nil, err
	}

	resp.CreatedAt = created_at.Format(time.RFC3339)
	resp.UpdatedAt = updated_at.Format(time.RFC3339)

	return resp, nil
}

func (s *PostAttachmentRepo) GetList(ctx context.Context, req *us.GetListAttachmentRequest) (*us.AttachmentList, error) {
	resp := &us.AttachmentList{}
	var (
		created_at, updated_at time.Time
	)
	offset := (req.Page - 1) * req.Limit

	var conditions []string
	if req.Search != "" {
		conditions = append(conditions, fmt.Sprintf(`post_id = '%v'`, req.Search))
	}

	var query string
	if len(conditions) > 0 {
		query = " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)

	rows, err := s.db.Query(ctx, `
		SELECT
			id,
			post_id,
			filepath,
			content_type,
			created_at,
			updated_at
		FROM post_attachment
	`+query)

	if err != nil {
		log.Println("error while getting all attachments:", err)
		return nil, err
	}

	defer rows.Close()

	var count int64

	for rows.Next() {
		var attachment us.Attachment
		count++
		err = rows.Scan(&attachment.Id, &attachment.PostId, &attachment.Filepath, &attachment.ContentType, &created_at, &updated_at)

		if err != nil {
			log.Println("error while scanning attachments:", err)
			return nil, err
		}
		attachment.CreatedAt = created_at.Format(time.RFC3339)
		attachment.UpdatedAt = updated_at.Format(time.RFC3339)

		resp.Items = append(resp.Items, &attachment)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

func (s *PostAttachmentRepo) Delete(ctx context.Context, req *us.AttachmentSingleRequest) (*emptypb.Empty, error) {
	_, err := s.db.Exec(ctx, `
		DELETE FROM post_attachment
		WHERE post_id = $1
	`, req.Id)

	if err != nil {
		log.Println("error while deleting attachment:", err)
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (s *PostAttachmentRepo) GetDefaultTags(ctx context.Context, req *us.GetDefaultTagsRequest) (*us.GetDefaultTagsResponse, error) {
	var response = &us.GetDefaultTagsResponse{}

	rows, err := s.db.Query(ctx, "SELECT name FROM default_tags")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		response.Tags = append(response.Tags, name)
	}

	return response, nil
}
