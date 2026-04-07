package postgres

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pgcontainer "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/Royal17x/subscription-service/internal/model"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	container, err := pgcontainer.Run(ctx,
		"postgres:16-alpine",
		pgcontainer.WithDatabase("testdb"),
		pgcontainer.WithUsername("postgres"),
		pgcontainer.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp"),
		),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, container.Terminate(ctx))
	})

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	m, err := migrate.New("file://../../../migrations", dsn)
	require.NoError(t, err)
	require.NoError(t, m.Up())

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func TestSubscriptionRepository_Create(t *testing.T) {
	pool := setupTestDB(t)
	repo := newSubscriptionRepository(pool)
	ctx := context.Background()

	userID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")
	startDate := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)

	sub := &model.Subscription{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      userID,
		StartDate:   startDate,
	}

	created, err := repo.Create(ctx, sub)
	require.NoError(t, err)
	assert.Greater(t, created.ID, int64(0))
	assert.Equal(t, "Yandex Plus", created.ServiceName)
	assert.Equal(t, 400, created.Price)
	assert.Equal(t, userID, created.UserID)
}

func TestSubscriptionRepository_GetByID(t *testing.T) {
	pool := setupTestDB(t)
	repo := newSubscriptionRepository(pool)
	ctx := context.Background()

	userID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")

	created, err := repo.Create(ctx, &model.Subscription{
		ServiceName: "Netflix",
		Price:       800,
		UserID:      userID,
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	t.Run("found", func(t *testing.T) {
		got, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, got.ID)
		assert.Equal(t, "Netflix", got.ServiceName)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByID(ctx, 999999)
		assert.ErrorIs(t, err, model.ErrNotFound)
	})
}

func TestSubscriptionRepository_Delete(t *testing.T) {
	pool := setupTestDB(t)
	repo := newSubscriptionRepository(pool)
	ctx := context.Background()

	userID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")

	created, err := repo.Create(ctx, &model.Subscription{
		ServiceName: "Spotify",
		Price:       200,
		UserID:      userID,
		StartDate:   time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		err := repo.Delete(ctx, created.ID)
		assert.NoError(t, err)

		_, err = repo.GetByID(ctx, created.ID)
		assert.ErrorIs(t, err, model.ErrNotFound)
	})

	t.Run("not found", func(t *testing.T) {
		err := repo.Delete(ctx, 999999)
		assert.ErrorIs(t, err, model.ErrNotFound)
	})
}

func TestSubscriptionRepository_TotalCost(t *testing.T) {
	pool := setupTestDB(t)
	repo := newSubscriptionRepository(pool)
	ctx := context.Background()

	userID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")

	_, err := repo.Create(ctx, &model.Subscription{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      userID,
		StartDate:   time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     timePtr(time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)),
	})
	require.NoError(t, err)

	_, err = repo.Create(ctx, &model.Subscription{
		ServiceName: "Netflix",
		Price:       800,
		UserID:      userID,
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     timePtr(time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)),
	})
	require.NoError(t, err)

	filter := model.TotalCostFilter{
		DateFrom: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
	}

	total, err := repo.TotalCost(ctx, filter)
	require.NoError(t, err)
	assert.Equal(t, int64(12000), total)
}

func timePtr(t time.Time) *time.Time {
	return &t
}
