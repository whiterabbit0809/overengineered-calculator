// internal/calculator/service.go
package calculator

import (
	"context"
	"errors"
	"fmt"

	"github.com/whiterabbit0809/overengineered-calculator/internal/history"
)

var ErrInvalidOperation = errors.New("invalid operation")
var ErrDivisionByZero = errors.New("division by zero")

type Service interface {
	Calculate(ctx context.Context, userID string, req CalculationRequest) (CalculationResult, error)
}

type service struct {
	historySvc history.Service
}

func NewService(historySvc history.Service) Service {
	return &service{historySvc: historySvc}
}

func (s *service) Calculate(ctx context.Context, userID string, req CalculationRequest) (CalculationResult, error) {
	// 1) Get previous result (state), default 0
	prevResult, err := s.historySvc.GetLatestResult(ctx, userID)
	if err != nil {
		return CalculationResult{}, err
	}

	// 2) Apply operation: result = prevResult (op) num
	newResult, expr, err := applyOperation(prevResult, req.Num, req.Operation)
	if err != nil {
		return CalculationResult{}, err
	}

	// 3) Save in history
	entry := &history.HistoryEntry{
		UserID:     userID,
		Expression: expr,
		Result:     newResult,
	}
	if err := s.historySvc.Record(ctx, entry); err != nil {
		return CalculationResult{}, err
	}

	return CalculationResult{
		Expression: expr,
		Result:     newResult,
	}, nil
}

func applyOperation(prev, num float64, op Operation) (float64, string, error) {
	switch op {
	case OpAdd:
		res := prev + num
		return res, fmt.Sprintf("%g + %g", prev, num), nil
	case OpSubtract:
		res := prev - num
		return res, fmt.Sprintf("%g - %g", prev, num), nil
	case OpMultiply:
		res := prev * num
		return res, fmt.Sprintf("%g * %g", prev, num), nil
	case OpDivide:
		if num == 0 {
			return 0, "", ErrDivisionByZero
		}
		res := prev / num
		return res, fmt.Sprintf("%g / %g", prev, num), nil
	default:
		return 0, "", ErrInvalidOperation
	}
}
