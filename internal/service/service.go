package service

import (
	"context"

	"github.com/Royal17x/subscription-service/internal/model"
)

type SubscriptionService interface {
	Create(ctx context.Context, req *model.SubscriptionRequest) (*model.SubscriptionResponse, error)
	GetByID(ctx context.Context, id int64) (*model.SubscriptionResponse, error)
	Update(ctx context.Context, id int64, req *model.SubscriptionRequest) (*model.SubscriptionResponse, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter model.SubscriptionFilter) ([]*model.SubscriptionResponse, error)
	TotalCost(ctx context.Context, filter model.TotalCostFilter) (int64, error)
}
