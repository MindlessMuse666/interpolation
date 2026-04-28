package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/user/interpolation/backend/core/interpolation"
)

func TestInterpolateEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/api/v1/interpolate", handleInterpolate)

	// Тестовые данные
	reqBody := InterpolateRequest{
		Method: "linear",
		Points: []interpolation.Point{
			{X: 0, Y: 0},
			{X: 2, Y: 4},
		},
		TargetX: 1,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/interpolate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp InterpolateResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 2.0, resp.Result)
	assert.Equal(t, "linear", resp.Method)
	assert.Len(t, resp.Curve, 100)
}
