package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mmomo/sezzle-calculator-project/backend/internal/calculator"
)

// CalculateHandler handles POST /api/calculate. It owns structural
// validation (400s); domain errors from calculator.Calculate become 422s.
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	op := calculator.Operation(req.Operation)
	if !calculator.IsValid(op) {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("unknown operation: %q", req.Operation))
		return
	}

	if req.A == nil {
		writeError(w, http.StatusBadRequest, `"a" is required and must be a number`)
		return
	}

	var b float64
	if !calculator.IsUnary(op) {
		if req.B == nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf(`"b" is required for operation %q`, req.Operation))
			return
		}
		b = *req.B
	}

	result, err := calculator.Calculate(op, *req.A, b)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, CalculateResponse{Result: result})
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
