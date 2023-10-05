package tests

import (
	"context"
	db "github.com/Bakhram74/amazon.git/db/sqlc"
	"github.com/Bakhram74/amazon.git/pkg/utils"

	"github.com/stretchr/testify/require"
	"testing"
)

func randomUser(t *testing.T) db.User {
	hashedPassword, err := utils.HashPassword(utils.RandomString(8))
	require.NoError(t, err)
	arg := db.CreateUserParams{
		Username:       utils.RandomString(6),
		PhoneNumber:    utils.RandomNumbers(9),
		HashedPassword: hashedPassword,
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.Role)
	require.Equal(t, false, user.IsBanned)
	require.Equal(t, "user", user.Role)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.PhoneNumber, user.PhoneNumber)
	return user
}

func TestCreateUser(t *testing.T) {
	randomUser(t)
}

func TestGetUser(t *testing.T) {
	user := randomUser(t)
	gotUser, err := testStore.GetUser(context.Background(), user.PhoneNumber)
	require.NoError(t, err)

	require.Equal(t, gotUser.CreatedAt, user.CreatedAt)
	require.Equal(t, gotUser.Username, user.Username)
	require.Equal(t, gotUser.HashedPassword, user.HashedPassword)
	require.Equal(t, gotUser.PhoneNumber, user.PhoneNumber)
}

func TestUpdateUser(t *testing.T) {
	user := randomUser(t)
	hashedPassword, _ := utils.HashPassword(utils.RandomString(8))
	params := db.PartialUpdateUserParams{
		ID:                user.ID,
		Username:          "Alex",
		UpdateUsername:    true,
		PhoneNumber:       utils.RandomNumbers(7),
		UpdatePhoneNumber: true,
		Password:          hashedPassword,
		UpdatePassword:    true,
	}
	updateUser, err := testStore.PartialUpdateUser(context.Background(), params)
	require.NoError(t, err)

	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.Role)
	require.Equal(t, false, user.IsBanned)

	require.NotEqual(t, user.Username, updateUser.Username)
	require.NotEqual(t, user.PhoneNumber, updateUser.PhoneNumber)
	require.NotEqual(t, user.HashedPassword, updateUser.HashedPassword)
	require.Equal(t, "Alex", updateUser.Username)
	require.Equal(t, 7, len(updateUser.PhoneNumber))
}

func TestUpdateUserName(t *testing.T) {
	user := randomUser(t)
	params := db.PartialUpdateUserParams{
		ID:             user.ID,
		Username:       "Pedro",
		UpdateUsername: true,
	}
	updateUser, err := testStore.PartialUpdateUser(context.Background(), params)
	require.NoError(t, err)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.Role)
	require.Equal(t, false, user.IsBanned)
	require.NotEqual(t, user.Username, updateUser.Username)
	require.Equal(t, user.PhoneNumber, updateUser.PhoneNumber)
	require.Equal(t, user.HashedPassword, updateUser.HashedPassword)
	require.Equal(t, "Pedro", updateUser.Username)
	require.Equal(t, user.PhoneNumber, updateUser.PhoneNumber)
}
func TestUpdateUserPhoneNumber(t *testing.T) {
	user := randomUser(t)
	params := db.PartialUpdateUserParams{
		ID:                user.ID,
		PhoneNumber:       utils.RandomNumbers(4),
		UpdatePhoneNumber: true,
	}
	updateUser, err := testStore.PartialUpdateUser(context.Background(), params)
	require.NoError(t, err)

	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.Role)
	require.Equal(t, false, user.IsBanned)

	require.Equal(t, user.Username, updateUser.Username)
	require.NotEqual(t, user.PhoneNumber, updateUser.PhoneNumber)
	require.Equal(t, 4, len(updateUser.PhoneNumber))
	require.Equal(t, user.HashedPassword, updateUser.HashedPassword)
	require.Equal(t, user.Username, updateUser.Username)
}

func TestUpdateUserPassword(t *testing.T) {
	hashedPassword, _ := utils.HashPassword(utils.RandomString(8))
	user := randomUser(t)
	params := db.PartialUpdateUserParams{
		ID:             user.ID,
		Password:       hashedPassword,
		UpdatePassword: true,
	}
	updateUser, err := testStore.PartialUpdateUser(context.Background(), params)
	require.NoError(t, err)

	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.Role)
	require.Equal(t, false, user.IsBanned)

	require.Equal(t, user.Username, updateUser.Username)
	require.NotEqual(t, user.HashedPassword, updateUser.HashedPassword)
	require.Equal(t, user.PhoneNumber, updateUser.PhoneNumber)
	require.Equal(t, user.Username, updateUser.Username)
}
