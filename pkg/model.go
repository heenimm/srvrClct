package pkg

type CalculationRequest struct {
	ID         int    `json:"id"`
	Expression string `json:"expression"`
}

type CalculationResponse struct {
	ID     string   `json:"id,omitempty"`
	Result *float64 `json:"result,omitempty"`
	Error  string   `json:"error,omitempty"`
}

type ExpressionRequest struct{}

type ExpressionResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Result string `json:"result,omitempty"`
}

type TaskRequest struct {
	ID            string `json:"id"`
	Arg1          string `json:"arg1"`
	Arg2          string `json:"arg2"`
	Operation     string `json:"operation"`
	OperationTime int    `json:"operation_time"`
}

type TaskResponse struct {
	ID     int     `json:"id"`
	Result float64 `json:"result"`
}
