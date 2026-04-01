package service

import (
	"context"
	"fmt"
	"github.com/Royal17x/subscription-service/internal/model"
	"github.com/Royal17x/subscription-service/internal/repository"
	"github.com/google/uuid"
	"time"
)

const dateTemplate = "01-2006"

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) Create(ctx context.Context, req *model.SubscriptionRequest) (*model.SubscriptionResponse, error) {
	sub, err := parseFromRequest(req)
	if err != nil {
		return nil, err
	}

	created, err := s.repo.Create(ctx, sub)
	if err != nil {
		return nil, fmt.Errorf("create subscription: %w", err)
	}
	return parseToResponse(created), nil
}

func (s *subscriptionService) GetByID(ctx context.Context, id int64) (*model.SubscriptionResponse, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseToResponse(sub), nil
}

func (s *subscriptionService) Update(ctx context.Context, id int64, req *model.SubscriptionRequest) (*model.SubscriptionResponse, error) {
	sub, err := parseFromRequest(req)
	if err != nil {
		return nil, err
	}
	sub.ID = id

	updated, err := s.repo.Update(ctx, sub)
	if err != nil {
		return nil, fmt.Errorf("update subscription: %w", err)
	}
	return parseToResponse(updated), nil
}

func (s *subscriptionService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *subscriptionService) List(ctx context.Context, filter model.SubscriptionFilter) ([]*model.SubscriptionResponse, error) {
	subs, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}

	result := make([]*model.SubscriptionResponse, 0, len(subs))
	for _, sub := range subs {
		result = append(result, parseToResponse(sub))
	}
	return result, nil
}

func (s *subscriptionService) TotalCost(ctx context.Context, filter model.TotalCostFilter) (int64, error) {
	total, err := s.repo.TotalCost(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("get total cost: %w", err)
	}
	return total, nil
}

func parseFromRequest(req *model.SubscriptionRequest) (*model.Subscription, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, model.ErrInvalidUUID
	}

	startDate, err := parseDate(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}

	sub := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
	}

	if req.EndDate != nil {
		endDate, err := parseDate(*req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date: %w", err)
		}
		if !endDate.After(startDate) {
			return nil, model.ErrInvalidDateRange
		}
		sub.EndDate = &endDate
	}
	return sub, nil
}

func parseToResponse(sub *model.Subscription) *model.SubscriptionResponse {
	resp := &model.SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID.String(),
		StartDate:   sub.StartDate.Format(dateTemplate),
		CreatedAt:   sub.CreatedAt.Format(time.RFC3339),
	}

	if sub.EndDate != nil {
		formatted := sub.EndDate.Format(dateTemplate)
		resp.EndDate = &formatted
	}
	return resp
}

func parseDate(date string) (time.Time, error) {
	t, err := time.Parse(dateTemplate, date)
	if err != nil {
		return time.Time{}, fmt.Errorf("expected format MM-YYYY, got %q", date)
	}
	return t, nil
}
