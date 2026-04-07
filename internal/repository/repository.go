package repository

import (
	"context"

	"github.com/Royal17x/subscription-service/internal/model"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *model.Subscription) (*model.Subscription, error)
	GetByID(ctx context.Context, id int64) (*model.Subscription, error)
	Update(ctx context.Context, sub *model.Subscription) (*model.Subscription, error)
	PartialUpdate(ctx context.Context, id int64, req *model.SubscriptionUpdateRequest) (*model.Subscription, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter model.SubscriptionFilter) ([]*model.Subscription, error)
	TotalCost(ctx context.Context, filter model.TotalCostFilter) (int64, error)
}
