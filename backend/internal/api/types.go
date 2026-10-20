package api

// CalculateRequest is the JSON body for POST /api/calculate.
//
// A and B are pointers so a missing field is distinguishable from an
// explicit 0, which the handler needs for "required" validation.
type CalculateRequest struct {
	Operation string   `json:"operation"`
	A         *float64 `json:"a"`
	B         *float64 `json:"b"`
}

// CalculateResponse is the JSON body for a successful calculation.
type CalculateResponse struct {
	Result float64 `json:"result"`
}

// ErrorResponse is the JSON body for a failed request (4xx).
type ErrorResponse struct {
	Error string `json:"error"`
}
