package tests

import (
	"context"
	db "github.com/Bakhram74/amazon-backend.git/db/sqlc"
	"github.com/Bakhram74/amazon-backend.git/pkg/utils"
	"github.com/stretchr/testify/require"

	"testing"
)

func randomProduct(t *testing.T, categoryId int32) db.Product {

	strArray := []string{"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSWdTM8I3dmUKDq1WzCMOsn85SqmqgmBXJ9Dg&usqp=CAU"}
	name := utils.RandomString(7)
	arg := db.CreateProductParams{
		CategoryID:  categoryId,
		Name:        name,
		Slug:        name,
		Description: utils.RandomString(30),
		Price:       777,
		Images:      strArray,
	}

	product, err := testStore.CreateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, product.CreatedAt)
	require.NotZero(t, product.ID)
	require.Equal(t, arg.Name, name)

	return product
}

func TestCreateProduct(t *testing.T) {
	category := randomCategory(t)
	product := randomProduct(t, category.ID)
	user := randomUser(t)
	randomReview(t, user.ID, product.ID)
}
