package history

import "context"

type Service interface {
	Record(ctx context.Context, entry *HistoryEntry) error
	List(ctx context.Context, userID string, limit, offset int) ([]HistoryEntry, error)
	GetLatestResult(ctx context.Context, userID string) (float64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Record(ctx context.Context, entry *HistoryEntry) error {
	return s.repo.Create(ctx, entry)
}

func (s *service) List(ctx context.Context, userID string, limit, offset int) ([]HistoryEntry, error) {
	return s.repo.ListByUser(ctx, userID, limit, offset)
}

func (s *service) GetLatestResult(ctx context.Context, userID string) (float64, error) {
	return s.repo.GetLatestResult(ctx, userID)
}
