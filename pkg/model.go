package pkg

type CalculationRequest struct {
	Expression string `json:"expression"`
}

type CalculationResponse struct {
	Result *float64 `json:"result,omitempty"`
	Error  string   `json:"error,omitempty"`
}
