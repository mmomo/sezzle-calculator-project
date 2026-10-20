package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantResult *float64
		wantErrSub string
	}{
		{
			name:       "add success",
			body:       `{"operation":"add","a":2,"b":3}`,
			wantStatus: http.StatusOK,
			wantResult: ptr(5),
		},
		{
			name:       "sqrt success without b",
			body:       `{"operation":"sqrt","a":16}`,
			wantStatus: http.StatusOK,
			wantResult: ptr(4),
		},
		{
			name:       "sqrt ignores extra b",
			body:       `{"operation":"sqrt","a":16,"b":99}`,
			wantStatus: http.StatusOK,
			wantResult: ptr(4),
		},
		{
			name:       "percentage success",
			body:       `{"operation":"percentage","a":50,"b":200}`,
			wantStatus: http.StatusOK,
			wantResult: ptr(100),
		},
		{
			name:       "malformed json",
			body:       `{"operation":`,
			wantStatus: http.StatusBadRequest,
			wantErrSub: "invalid request body",
		},
		{
			name:       "unknown operation",
			body:       `{"operation":"modulo","a":1,"b":2}`,
			wantStatus: http.StatusBadRequest,
			wantErrSub: "unknown operation",
		},
		{
			name:       "missing a",
			body:       `{"operation":"add","b":2}`,
			wantStatus: http.StatusBadRequest,
			wantErrSub: `"a" is required`,
		},
		{
			name:       "missing b for binary operation",
			body:       `{"operation":"add","a":1}`,
			wantStatus: http.StatusBadRequest,
			wantErrSub: `"b" is required`,
		},
		{
			name:       "divide by zero",
			body:       `{"operation":"divide","a":1,"b":0}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantErrSub: "division by zero",
		},
		{
			name:       "negative sqrt",
			body:       `{"operation":"sqrt","a":-4}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantErrSub: "square root",
		},
		{
			name:       "exponentiate overflow",
			body:       `{"operation":"exponentiate","a":10,"b":400}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantErrSub: "not a finite number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			CalculateHandler(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}

			if tt.wantResult != nil {
				var resp CalculateResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp.Result != *tt.wantResult {
					t.Fatalf("result = %v, want %v", resp.Result, *tt.wantResult)
				}
			}

			if tt.wantErrSub != "" {
				var resp ErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Fatalf("decode error response: %v", err)
				}
				if !strings.Contains(resp.Error, tt.wantErrSub) {
					t.Fatalf("error = %q, want substring %q", resp.Error, tt.wantErrSub)
				}
			}
		})
	}
}

func ptr(f float64) *float64 {
	return &f
}
