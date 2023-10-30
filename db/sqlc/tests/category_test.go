package tests

import (
	"context"
	db "github.com/Bakhram74/amazon-backend.git/db/sqlc"
	"github.com/Bakhram74/amazon-backend.git/pkg/utils"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func randomCategory(t *testing.T) db.Category {
	id, err := strconv.Atoi(utils.RandomNumbers(10))
	require.NoError(t, err)
	categoryArgs := db.CreateCategoryParams{
		ID:   int32(id),
		Name: utils.RandomString(10),
	}
	category, err := testStore.CreateCategory(context.Background(), categoryArgs)
	require.NoError(t, err)

	require.NotZero(t, category.CreatedAt)
	require.NotZero(t, category.ID)

	require.Equal(t, categoryArgs.Name, category.Name)
	require.Equal(t, categoryArgs.ID, category.ID)
	return category
}
func TestCreateCategory(t *testing.T) {
	randomCategory(t)
}
