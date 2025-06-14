package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Thingsly/backend/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestHomepageHandler(t *testing.T) {
	userApi := api.UserApi{}
	r := SetUpRouter()
	r.GET("/api/v1/user", userApi.HandleUserListByPage)

	req, _ := http.NewRequest("GET", "/api/v1/user?page=1&page_size=10", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-token",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InI1NjYzMzU5LTU0ZjMtZTIzYi04MDFiLWUzYjgyNDQ2Mzc4ZSIsImVtYWlsIjoic3VwZXJAc3VwZXIuY24iLCJjcmVhdGVfdGltZSI6IjIwMjQtMDMtMDZUMTQ6NDE6MTEuNzE1MDYwNCswODowMCIsImF1dGhvcml0eSI6IlNZU19BRE1JTiIsInRlbmFudF9pZCI6ImFhYWFhYWFhIn0.q2WI_eQ0837jAqCkE-Tj27IZ5C7qpYCaP9lJm1qVF2k",
	)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)

	assert.Equal(t, http.StatusOK, w.Code)
	t.Log("responseData:" + string(responseData))
}
