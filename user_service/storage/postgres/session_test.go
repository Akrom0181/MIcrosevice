package postgres_test

import (
	"context"
	"testing"
	"time"
	"user_service/genproto/user_service"
	"user_service/storage/postgres"

	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestSessionRepo_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := postgres.NewSessionRepo(db)
	ctx := context.Background()

	req := &user_service.Session{
		UserId:       "9e129b9e-795e-4942-9d7d-639ccc92953d",
		IpAddress:    "127.0.0.1",
		Platform:     "web",
		IsActive:     true,
		UserAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36",
		LastActiveAt: time.Now().UTC().Format(time.RFC3339),
		ExpiresAt:    time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339),
	}

	session, err := repo.Create(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, session)

	assert.Equal(t, req.UserId, session.UserId)
	assert.Equal(t, req.IpAddress, session.IpAddress)
	assert.Equal(t, req.Platform, session.Platform)
	assert.Equal(t, req.IsActive, session.IsActive)
	assert.Equal(t, req.UserAgent, session.UserAgent)

	expectedLastActiveAt, _ := time.Parse(time.RFC3339, req.LastActiveAt)
	actualLastActiveAt, _ := time.Parse(time.RFC3339, session.LastActiveAt)
	require.WithinDuration(t, expectedLastActiveAt, actualLastActiveAt, time.Second)

	expectedExpiresAt, _ := time.Parse(time.RFC3339, req.ExpiresAt)
	actualExpiresAt, _ := time.Parse(time.RFC3339, session.ExpiresAt)
	require.WithinDuration(t, expectedExpiresAt, actualExpiresAt, time.Second)
}

func TestSessionRepo_GetSingle(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := postgres.NewSessionRepo(db)

	ctx := context.Background()
	req := &user_service.SessionSingleRequest{
		Id: "43760ef9-949d-4e67-ba61-6192f40ce5f1",
	}

	session, err := repo.GetSingle(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, "43760ef9-949d-4e67-ba61-6192f40ce5f1", session.Id)
	assert.Equal(t, "127.0.0.1", session.IpAddress)
	assert.Equal(t, "web", session.Platform)
	assert.Equal(t, true, session.IsActive)
	assert.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36", session.UserAgent)

}

func TestSessionRepo_GetList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := postgres.NewSessionRepo(db)

	ctx := context.Background()
	req := &user_service.GetListSessionRequest{
		Limit:  1,
		Page:   1,
		Search: "9e129b9e-795e-4942-9d7d-639ccc92953d",
	}

	sessions, err := repo.GetList(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, "43760ef9-949d-4e67-ba61-6192f40ce5f1", sessions.Sessions[0].Id)
	assert.Equal(t, "::1", sessions.Sessions[0].IpAddress)
	assert.Equal(t, "web", sessions.Sessions[0].Platform)
	assert.Equal(t, true, sessions.Sessions[0].IsActive)
	assert.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36", sessions.Sessions[0].UserAgent)
}

func TestSessionRepo_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := postgres.NewSessionRepo(db)

	ctx := context.Background()
	req := &user_service.SessionSingleRequest{
		Id: "95c5f6bd-8840-4893-ae6f-d97ee3a6d6ab",
	}

	_, err := repo.Delete(ctx, req)
	require.NoError(t, err)
}

func TestSessionRepo_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := postgres.NewSessionRepo(db)

	ctx := context.Background()
	req := &user_service.Session{
		Id:           "43760ef9-949d-4e67-ba61-6192f40ce5f1",
		IpAddress:    "127.0.0.1",
		LastActiveAt: "2022-12-12T12:12:12+00:00",
		IsActive:     true,
	}

	_, err := repo.Update(ctx, req)
	require.NoError(t, err)
}
