package tests

import (
	"context"
	db "github.com/Bakhram74/amazon-backend.git/db/sqlc"
	"github.com/Bakhram74/amazon-backend.git/pkg/utils"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func randomReview(t *testing.T, userId int64, productId int32) db.Review {
	rating, err := strconv.Atoi(utils.RandomNumbers(1))

	require.NoError(t, err)
	reviewArgs := db.CreateReviewParams{
		UserID:    userId,
		ProductID: productId,
		Rating:    int32(rating),
		Text:      utils.RandomString(20),
	}
	review, err := testStore.CreateReview(context.Background(), reviewArgs)
	require.NoError(t, err)

	require.NotZero(t, review.CreatedAt)
	require.NotZero(t, review.ID)

	require.Equal(t, review.Text, reviewArgs.Text)
	require.Equal(t, review.Rating, reviewArgs.Rating)
	require.Equal(t, review.ProductID, reviewArgs.ProductID)
	require.Equal(t, review.UserID, reviewArgs.UserID)
	return review
}
