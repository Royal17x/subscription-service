package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Royal17x/subscription-service/internal/model"
	"github.com/Royal17x/subscription-service/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type subscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) repository.SubscriptionRepository {
	return &subscriptionRepository{pool: pool}
}

func (r *subscriptionRepository) Create(ctx context.Context, sub *model.Subscription) (*model.Subscription, error) {
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at`

	row := r.pool.QueryRow(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	)
	return scanSubscription(row)
}

func (r *subscriptionRepository) GetByID(ctx context.Context, id int64) (*model.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at
		FROM subscriptions
		WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)
	sub, err := scanSubscription(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return sub, nil
}

func (r *subscriptionRepository) Update(ctx context.Context, sub *model.Subscription) (*model.Subscription, error) {
	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
		WHERE id = $6
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at`

	row := r.pool.QueryRow(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)
	updatedSub, err := scanSubscription(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return updatedSub, nil
}

func (r *subscriptionRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM subscriptions
		WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}
	if result.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *subscriptionRepository) List(ctx context.Context, filter model.SubscriptionFilter) ([]*model.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at
		FROM subscriptions
		WHERE ($1::uuid IS NULL OR user_id = $1)
			AND ($2::varchar IS NULL OR service_name = $2)
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, filter.UserID, filter.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []*model.Subscription
	for rows.Next() {
		sub, err := scanSubscription(rows)
		if err != nil {
			return nil, fmt.Errorf("list subscriptions: %w", err)
		}
		subs = append(subs, sub)
	}
	return subs, rows.Err()
}

func (r *subscriptionRepository) TotalCost(ctx context.Context, filter model.TotalCostFilter) (int64, error) {
	query := `
		SELECT COALESCE(SUM(price),0)
		FROM subscriptions
		WHERE ($1::uuid is NULL OR user_id = $1)
			AND ($2::varchar is NULL OR service_name = $2)
			AND start_date <= $4
			AND (end_date IS NULL OR end_date >= $3)`

	var totalCost int64
	err := r.pool.QueryRow(ctx, query,
		filter.UserID,
		filter.ServiceName,
		filter.DateFrom,
		filter.DateTo,
	).Scan(&totalCost)
	if err != nil {
		return 0, fmt.Errorf("total cost: %w", err)
	}
	return totalCost, nil
}

type scanner interface {
	Scan(row ...any) error
}

func scanSubscription(s scanner) (*model.Subscription, error) {
	var sub model.Subscription
	err := s.Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}
