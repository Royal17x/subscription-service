package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Royal17x/subscription-service/internal/mocks"
	"github.com/Royal17x/subscription-service/internal/model"
)

func TestSubscriptionService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockSubscriptionRepository(ctrl)
	svc := NewSubscriptionService(repo)

	userID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")
	startDate := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)

	t.Run("success", func(t *testing.T) {
		req := &model.SubscriptionRequest{
			ServiceName: "Spotify",
			Price:       400,
			UserID:      userID.String(),
			StartDate:   "07-2025",
		}

		expected := &model.Subscription{
			ID:          1,
			ServiceName: "Spotify",
			Price:       400,
			UserID:      userID,
			StartDate:   startDate,
			CreatedAt:   time.Now(),
		}

		repo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(expected, nil)

		resp, err := svc.Create(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, int64(1), resp.ID)
		assert.Equal(t, "Spotify", resp.ServiceName)
		assert.Equal(t, "07-2025", resp.StartDate)
		assert.Equal(t, 400, resp.Price)
	})

	t.Run("invalid uuid", func(t *testing.T) {
		req := &model.SubscriptionRequest{
			ServiceName: "Spotify",
			Price:       400,
			UserID:      "not-a-uuid",
			StartDate:   "07-2025",
		}

		_, err := svc.Create(context.Background(), req)
		assert.Error(t, err, model.ErrInvalidUUID)
	})

	t.Run("invalid start date", func(t *testing.T) {
		req := &model.SubscriptionRequest{
			ServiceName: "Spotify",
			Price:       400,
			UserID:      userID.String(),
			StartDate:   "2025-07",
		}

		_, err := svc.Create(context.Background(), req)
		assert.Error(t, err)
	})

	t.Run("end_date before start_date", func(t *testing.T) {
		endDate := "06-2025"
		req := &model.SubscriptionRequest{
			ServiceName: "Spotify",
			Price:       400,
			UserID:      userID.String(),
			StartDate:   "07-2025",
			EndDate:     &endDate,
		}

		_, err := svc.Create(context.Background(), req)
		assert.ErrorIs(t, err, model.ErrInvalidDateRange)
	})
}

func TestSubscriptionService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockSubscriptionRepository(ctrl)
	svc := NewSubscriptionService(repo)

	userID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")

	t.Run("success", func(t *testing.T) {
		expected := &model.Subscription{
			ID:          1,
			ServiceName: "Netflix",
			Price:       800,
			UserID:      userID,
			StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			CreatedAt:   time.Now(),
		}

		repo.EXPECT().
			GetByID(gomock.Any(), int64(1)).
			Return(expected, nil)

		resp, err := svc.GetByID(context.Background(), 1)
		require.NoError(t, err)
		assert.Equal(t, int64(1), resp.ID)
		assert.Equal(t, "Netflix", resp.ServiceName)
	})

	t.Run("not found", func(t *testing.T) {
		repo.EXPECT().
			GetByID(gomock.Any(), int64(999)).
			Return(nil, model.ErrNotFound)

		_, err := svc.GetByID(context.Background(), 999)
		assert.ErrorIs(t, err, model.ErrNotFound)
	})
}

func TestSubscriptionService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockSubscriptionRepository(ctrl)
	svc := NewSubscriptionService(repo)

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().
			Delete(gomock.Any(), int64(1)).
			Return(nil)

		err := svc.Delete(context.Background(), 1)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		repo.EXPECT().
			Delete(gomock.Any(), int64(999)).
			Return(model.ErrNotFound)

		err := svc.Delete(context.Background(), 999)
		assert.ErrorIs(t, err, model.ErrNotFound)
	})
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:  "valid date",
			input: "07-2025",
			want:  time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:    "wrong format",
			input:   "2025-07",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDate(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

