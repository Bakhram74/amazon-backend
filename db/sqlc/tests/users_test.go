package tests

import (
	"context"
	db "github.com/Bakhram74/amazon-backend.git/db/sqlc"
	"github.com/Bakhram74/amazon-backend.git/pkg/utils"

	"github.com/stretchr/testify/require"
	"testing"
)

func randomUser(t *testing.T) db.User {
	hashedPassword, err := utils.HashPassword(utils.RandomString(8))
	require.NoError(t, err)
	arg := db.CreateUserParams{
		Name:           utils.RandomString(6),
		Email:          utils.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.ID)
	require.NotEmpty(t, user.AvatarPath)
	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	return user
}

func TestCreateUser(t *testing.T) {
	randomUser(t)
}
