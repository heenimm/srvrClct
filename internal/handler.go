package internal

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"serverCalc/pkg"
)

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	var response pkg.CalculationResponse

	if r.Method != http.MethodPost {
		response.Error = "method  must be POST 'http://localhost:8082/api/v1/calculate'"
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request pkg.CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error = "Expression is not valid"
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	result, err := pkg.Calculate(request.Expression)

	if err != nil && errors.Is(err, pkg.ErrInternalError) {
		response.Error = "Internal server error"
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err != nil {
		response.Error = "Expression is not valid"
		w.WriteHeader(http.StatusUnprocessableEntity)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			panic(err)
		}
		return
	}

	response.Result = &result
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Panicln("Error encoding response:", err)
	}
}
