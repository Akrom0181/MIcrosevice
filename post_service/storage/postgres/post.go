package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	us "post_service/genproto/post_service"
	"post_service/storage"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PostRepo struct {
	db *pgxpool.Pool
}

func NewPostRepo(db *pgxpool.Pool) storage.PostRepoI {
	return &PostRepo{
		db: db,
	}
}

// Create implements storage.PostRepoI.
func (s *PostRepo) Create(ctx context.Context, req *us.Post) (*us.Post, error) {
	id := uuid.NewString()

	_, err := s.db.Exec(ctx, `
		INSERT INTO posts (
			id,
			owner_id,
			tags,
			content,
			status
		) VALUES (
			$1, $2, $3, $4, $5
		)`, id, req.OwnerId, req.Tags, req.Content, req.Status)

	if err != nil {
		log.Println("error while creating post in storage", err)
		return nil, err
	}

	post, err := s.GetSingle(ctx, &us.PostSingleRequest{Id: id})
	if err != nil {
		log.Println("error while getting post by id after creating", err)
		return nil, err
	}
	return post, nil
}

// GetByID implements storage.PostRepoI.
func (s *PostRepo) GetSingle(ctx context.Context, req *us.PostSingleRequest) (*us.Post, error) {
	resp := &us.Post{}

	var (
		created_at, updated_at time.Time
	)

	tags := []byte{}
	err := s.db.QueryRow(ctx, `
	        SELECT 
			id,
			owner_id,
			tags,
			content,
			status,
	        created_at,
	        updated_at
	        FROM posts 
	    WHERE id=$1`, req.Id).Scan(&resp.Id, &resp.OwnerId, &tags, &resp.Content, &resp.Status, &created_at, &updated_at)

	if err != nil {
		log.Println("error while getting post by id", err)
		return &us.Post{}, err
	}

	err = json.Unmarshal(tags, &resp.Tags)
	if err != nil {
		return &us.Post{}, err
	}

	resp.CreatedAt = created_at.Format(time.RFC3339)
	resp.UpdatedAt = updated_at.Format(time.RFC3339)

	return resp, nil
}

// GetList implements storage.PostRepoI.
func (s *PostRepo) GetList(ctx context.Context, req *us.GetListPostRequest) (*us.PostList, error) {
	resp := &us.PostList{}

	var (
		createdAt, updatedAt time.Time
	)

	offset := (req.Page - 1) * req.Limit

	baseQuery := `
    SELECT 
        p.id, 
        p.owner_id,
		p.tags, 
        p.content, 
        p.status, 
        p.created_at, 
        p.updated_at,
        COALESCE(json_agg(DISTINCT pa_json), '[]'::json) AS attachments
    FROM posts p
    LEFT JOIN (
        SELECT pa.post_id, jsonb_build_object(
            'id', pa.id,
            'filepath', pa.filepath,
            'content_type', pa.content_type,
            'created_at', pa.created_at,
            'updated_at', pa.updated_at
        ) AS pa_json
        FROM post_attachment pa
    ) pa ON pa.post_id = p.id
    `

	whereClause := ""
	if req.Search != "" {
		whereClause = fmt.Sprintf(" WHERE p.owner_id = '%s'", req.Search)
	}

	tailClause := fmt.Sprintf(`
    %s
    GROUP BY p.id
    ORDER BY p.created_at DESC
    OFFSET %v LIMIT %v`, whereClause, offset, req.Limit)

	query := baseQuery + tailClause

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		log.Println("Error while getting posts:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post us.Post
		var attachmentsJSON string

		err = rows.Scan(&post.Id, &post.OwnerId, &post.Tags, &post.Content, &post.Status, &createdAt, &updatedAt, &attachmentsJSON)
		if err != nil {
			log.Println("error while scanning posts:", err)
			return nil, err
		}

		post.CreatedAt = createdAt.Format(time.RFC3339)
		post.UpdatedAt = updatedAt.Format(time.RFC3339)

		if err := json.Unmarshal([]byte(attachmentsJSON), &post.Attachments); err != nil {
			log.Println("error unmarshalling attachments:", err)
			return nil, err
		}

		resp.Items = append(resp.Items, &post)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = int64(len(resp.Items))
	return resp, nil
}

// Update implements storage.PostRepoI.
func (s *PostRepo) Update(ctx context.Context, req *us.Post) (*us.Post, error) {

	_, err := s.db.Exec(ctx, `
        UPDATE posts SET
		    content=$1,
		    status=$2,
            updated_at = NOW()
        WHERE id = $3`, req.Content, req.Status, req.Id)

	if err != nil {
		log.Println("error while updating post in storage", err)
		return nil, err
	}

	post, err := s.GetSingle(ctx, &us.PostSingleRequest{Id: req.Id})
	if err != nil {
		log.Println("error while getting updated post by id", err)
		return nil, err
	}

	return post, nil
}

// Delete implements storage.PostRepoI.
func (s *PostRepo) Delete(ctx context.Context, req *us.PostSingleRequest) (*emptypb.Empty, error) {
	_, err := s.db.Exec(ctx, `
		DELETE FROM posts
		WHERE id = $1
	`, req.Id)

	if err != nil {
		log.Println("error while deleting post:", err)
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
