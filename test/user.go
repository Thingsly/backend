package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HustIoTPlatform/backend/internal/api"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestHomepageHandler(t *testing.T) {
	userApi := api.UserApi{}
	// mockResponse := `{"message":"Welcome to the Tech Company listing API with Golang"}`
	r := SetUpRouter()
	r.GET("/api/v1/user", userApi.HandleUserListByPage)
	req, _ := http.NewRequest("GET", "/api/v1/user?page=1&page_size=10", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InI1NjYzMzU5LTU0ZjMtZTIzYi04MDFiLWUzYjgyNDQ2Mzc4ZSIsImVtYWlsIjoic3VwZXJAc3VwZXIuY24iLCJjcmVhdGVfdGltZSI6IjIwMjQtMDMtMDZUMTQ6NDE6MTEuNzE1MDYwNCswODowMCIsImF1dGhvcml0eSI6IlNZU19BRE1JTiIsInRlbmFudF9pZCI6ImFhYWFhYWFhIn0.q2WI_eQ0837jAqCkE-Tj27IZ5C7qpYCaP9lJm1qVF2k",
	)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)

	assert.Equal(t, http.StatusOK, w.Code)
	// assert.Equal(t, mockResponse, string(responseData))  
	t.Log("resoponseData:" + string(responseData))
	t.Log("resoponseData:" + string(responseData))
}
