package calculator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/whiterabbit0809/overengineered-calculator/internal/history"
)

//
// Test fakes
//

// fakeHistoryService implements history.Service for tests.
// Only GetLatestResult and Record are actually used by the calculator,
// but List is added so this type fully satisfies history.Service.
type fakeHistoryService struct {
	latestResult float64
	latestErr    error

	recordedEntries []*history.HistoryEntry
	recordErr       error
}

func (f *fakeHistoryService) GetLatestResult(ctx context.Context, userID string) (float64, error) {
	return f.latestResult, f.latestErr
}

func (f *fakeHistoryService) Record(ctx context.Context, entry *history.HistoryEntry) error {
	f.recordedEntries = append(f.recordedEntries, entry)
	return f.recordErr
}

func (f *fakeHistoryService) List(ctx context.Context, userID string, limit, offset int) ([]history.HistoryEntry, error) {
	// not used in calculator tests
	return nil, nil
}

// helper to build the concrete *service under test
func newTestCalcServiceWithHistory(hs history.Service) *service {
	return NewService(hs).(*service)
}

//
// Tests
//

// Test ADD from initial state (no history → previous result = 0).
func TestCalculate_AddFromZero(t *testing.T) {
	fh := &fakeHistoryService{
		latestResult: 0, // simulate "no previous result" case handled by history
	}
	svc := newTestCalcServiceWithHistory(fh)

	req := CalculationRequest{
		Num:       1,
		Operation: OpAdd,
	}

	res, err := svc.Calculate(context.Background(), "user-123", req)
	require.NoError(t, err)

	// 0 + 1 = 1
	assert.Equal(t, 1.0, res.Result)
	assert.Equal(t, "0 + 1", res.Expression)

	// A history entry must have been recorded
	require.Len(t, fh.recordedEntries, 1)
	entry := fh.recordedEntries[0]
	assert.Equal(t, "user-123", entry.UserID)
	assert.Equal(t, res.Expression, entry.Expression)
	assert.Equal(t, res.Result, entry.Result)
}

// Test SUBTRACT when there is a previous result.
func TestCalculate_SubtractWithPreviousResult(t *testing.T) {
	// Pretend last result was 10
	fh := &fakeHistoryService{
		latestResult: 10,
	}
	svc := newTestCalcServiceWithHistory(fh)

	req := CalculationRequest{
		Num:       3,
		Operation: OpSubtract,
	}

	res, err := svc.Calculate(context.Background(), "user-123", req)
	require.NoError(t, err)

	// 10 - 3 = 7
	assert.Equal(t, 7.0, res.Result)
	assert.Equal(t, "10 - 3", res.Expression)

	require.Len(t, fh.recordedEntries, 1)
	entry := fh.recordedEntries[0]
	assert.Equal(t, "user-123", entry.UserID)
	assert.Equal(t, res.Expression, entry.Expression)
	assert.Equal(t, res.Result, entry.Result)
}

// Test DIVIDE by zero: should return ErrDivisionByZero and not record history.
func TestCalculate_DivideByZero(t *testing.T) {
	fh := &fakeHistoryService{
		latestResult: 5,
	}
	svc := newTestCalcServiceWithHistory(fh)

	req := CalculationRequest{
		Num:       0,
		Operation: OpDivide,
	}

	_, err := svc.Calculate(context.Background(), "user-123", req)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrDivisionByZero)

	// No history entry must be recorded on failure
	assert.Len(t, fh.recordedEntries, 0)
}

// Test invalid operation: should return ErrInvalidOperation and not record.
func TestCalculate_InvalidOperation(t *testing.T) {
	fh := &fakeHistoryService{
		latestResult: 0,
	}
	svc := newTestCalcServiceWithHistory(fh)

	req := CalculationRequest{
		Num:       2,
		Operation: Operation("BOGUS"),
	}

	_, err := svc.Calculate(context.Background(), "user-123", req)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidOperation)

	assert.Len(t, fh.recordedEntries, 0)
}
