package postgres_test

import (
	"context"
	"testing"
	"user_service/genproto/user_service"
	"user_service/storage/postgres"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	db, err := pgxpool.Connect(context.Background(), "postgresql://akromjonotaboyev:1@localhost:5432/microservice")
	require.NoError(t, err)
	return db
}

func TestUserRepo_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := postgres.NewUserRepo(db)

	ctx := context.Background()
	req := &user_service.User{
		UserType: "admin",
		UserRole: "user",
		FullName: "John Doe",
		UserName: "johndoe",
		Email:    "johndoe@example.com",
		Password: "password123",
		Gender:   "male",
		Status:   "active",
	}

	createdUser, err := repo.Create(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, req.UserType, createdUser.UserType)
	assert.Equal(t, req.UserRole, createdUser.UserRole)
	assert.Equal(t, req.FullName, createdUser.FullName)
	assert.Equal(t, req.UserName, createdUser.UserName)
	assert.Equal(t, req.Email, createdUser.Email)
	assert.Equal(t, req.Password, createdUser.Password)
	assert.Equal(t, req.Gender, createdUser.Gender)
	assert.Equal(t, req.Status, createdUser.Status)
	assert.NotEmpty(t, createdUser.Id)
}

func TestUserRepo_Create_Error(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := postgres.NewUserRepo(db)

	ctx := context.Background()
	req := &user_service.User{
		UserType: "admin",
		UserRole: "user",
		FullName: "John Doe",
		UserName: "john_doe",
		Email:    "johndoe@gmails.com",
		Password: "password123",
		Gender:   "male",
		Status:   "active",
	}

	// Insert the user to create a conflict
	_, err := repo.Create(ctx, req)
	require.NoError(t, err)

	// Try to create the same user again to trigger an error
	_, err = repo.Create(ctx, req)
	assert.Error(t, err)
}

func TestUserRepo_GetSingle(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := postgres.NewUserRepo(db)

	ctx := context.Background()
	req := &user_service.UserSingleRequest{
		Id: "83c49d98-04bf-47ac-ba80-c1064bd870cf",
	}

	user, err := repo.GetSingle(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "83c49d98-04bf-47ac-ba80-c1064bd870cf", user.Id)
	assert.Equal(t, "admin", user.UserType)
	assert.Equal(t, "user", user.UserRole)
	assert.Equal(t, "John Doe", user.FullName)
	assert.Equal(t, "johndoe", user.UserName)
	assert.Equal(t, "johndoe@example.com", user.Email)
	assert.Equal(t, "password123", user.Password)
	assert.Equal(t, "male", user.Gender)
	assert.Equal(t, "active", user.Status)
}

func TestUserRepo(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := postgres.NewUserRepo(db)

	ctx := context.Background()
	req := &user_service.GetListUserRequest{
		Limit:  1,
		Page:   1,
		Search: "83c49d98-04bf-47ac-ba80-c1064bd870cf",
	}

	users, err := repo.GetList(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, "83c49d98-04bf-47ac-ba80-c1064bd870cf", users.Users[0].Id)
	assert.Equal(t, "admin", users.Users[0].UserType)
	assert.Equal(t, "user", users.Users[0].UserRole)
	assert.Equal(t, "John Doe", users.Users[0].FullName)
	assert.Equal(t, "johndoe", users.Users[0].UserName)
	assert.Equal(t, "johndoe@example.com", users.Users[0].Email)
	assert.Equal(t, "password123", users.Users[0].Password)
	assert.Equal(t, "male", users.Users[0].Gender)
	assert.Equal(t, "active", users.Users[0].Status)
}

func TestUserRepo_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := postgres.NewUserRepo(db)

	ctx := context.Background()
	req := &user_service.UserPrimaryKey{
		Id: "83c49d98-04bf-47ac-ba80-c1064bd870cf",
	}

	_, err := repo.Delete(ctx, req)
	require.NoError(t, err)
}
