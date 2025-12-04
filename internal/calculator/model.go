package calculator

type Operation string

const (
	OpAdd      Operation = "ADD"
	OpSubtract Operation = "SUBTRACT"
	OpMultiply Operation = "MULTIPLY"
	OpDivide   Operation = "DIVIDE"
)

type CalculationRequest struct {
	Num       float64   `json:"num"`
	Operation Operation `json:"operation"`
}

type CalculationResult struct {
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}
