package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Thingsly/backend/internal/api"
	"github.com/Thingsly/backend/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestCreateDevice(t *testing.T) {
	deviceApi := api.DeviceApi{}
	r := SetUpRouter()
	r.POST("/api/v1/device", deviceApi.CreateDevice)

	// Create request body
	name := "test-create-device"
	label := ""
	deviceConfigId := ""
	accessWay := "A"
	reqBody := model.CreateDeviceReq{
		Name:           &name,
		Label:          &label,
		DeviceConfigId: &deviceConfigId,
		AccessWay:      &accessWay,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/v1/device", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-token",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InI1NjYzMzU5LTU0ZjMtZTIzYi04MDFiLWUzYjgyNDQ2Mzc4ZSIsImVtYWlsIjoic3VwZXJAc3VwZXIuY24iLCJjcmVhdGVfdGltZSI6IjIwMjQtMDMtMDZUMTQ6NDE6MTEuNzE1MDYwNCswODowMCIsImF1dGhvcml0eSI6IlNZU19BRE1JTiIsInRlbmFudF9pZCI6ImFhYWFhYWFhIn0.q2WI_eQ0837jAqCkE-Tj27IZ5C7qpYCaP9lJm1qVF2k",
	)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)

	assert.Equal(t, http.StatusOK, w.Code)
	t.Log("Create device response:" + string(responseData))

	// Parse response to verify device was created
	var response map[string]interface{}
	err := json.Unmarshal(responseData, &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])

	// Verify device data
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "test-create-device", data["name"])
	assert.Equal(t, "A", data["access_way"])
}

func TestGetDeviceList(t *testing.T) {
	deviceApi := api.DeviceApi{}
	r := SetUpRouter()
	r.GET("/api/v1/device", deviceApi.HandleDeviceListByPage)

	req, _ := http.NewRequest("GET", "/api/v1/device?page=1&page_size=10", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-token",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InI1NjYzMzU5LTU0ZjMtZTIzYi04MDFiLWUzYjgyNDQ2Mzc4ZSIsImVtYWlsIjoic3VwZXJAc3VwZXIuY24iLCJjcmVhdGVfdGltZSI6IjIwMjQtMDMtMDZUMTQ6NDE6MTEuNzE1MDYwNCswODowMCIsImF1dGhvcml0eSI6IlNZU19BRE1JTiIsInRlbmFudF9pZCI6ImFhYWFhYWFhIn0.q2WI_eQ0837jAqCkE-Tj27IZ5C7qpYCaP9lJm1qVF2k",
	)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)

	assert.Equal(t, http.StatusOK, w.Code)
	t.Log("Get device list response:" + string(responseData))

	// Parse response to verify device list
	var response map[string]interface{}
	err := json.Unmarshal(responseData, &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])

	// Verify response has data field
	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["list"])
	assert.NotNil(t, data["total"])
}

func TestGetDeviceById(t *testing.T) {
	deviceApi := api.DeviceApi{}
	r := SetUpRouter()
	r.GET("/api/v1/device/detail/:id", deviceApi.HandleDeviceByID)

	// First create a device to get its ID
	name := "test-get-device"
	label := ""
	deviceConfigId := ""
	accessWay := "A"
	createReqBody := model.CreateDeviceReq{
		Name:           &name,
		Label:          &label,
		DeviceConfigId: &deviceConfigId,
		AccessWay:      &accessWay,
	}
	createJsonBody, _ := json.Marshal(createReqBody)

	createReq, _ := http.NewRequest("POST", "/api/v1/device", bytes.NewBuffer(createJsonBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("x-token",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InI1NjYzMzU5LTU0ZjMtZTIzYi04MDFiLWUzYjgyNDQ2Mzc4ZSIsImVtYWlsIjoic3VwZXJAc3VwZXIuY24iLCJjcmVhdGVfdGltZSI6IjIwMjQtMDMtMDZUMTQ6NDE6MTEuNzE1MDYwNCswODowMCIsImF1dGhvcml0eSI6IlNZU19BRE1JTiIsInRlbmFudF9pZCI6ImFhYWFhYWFhIn0.q2WI_eQ0837jAqCkE-Tj27IZ5C7qpYCaP9lJm1qVF2k",
	)

	createW := httptest.NewRecorder()
	r.ServeHTTP(createW, createReq)

	createResponseData, _ := io.ReadAll(createW.Body)
	var createResponse map[string]interface{}
	json.Unmarshal(createResponseData, &createResponse)
	deviceId := createResponse["data"].(map[string]interface{})["id"].(string)

	// Now get the device by ID
	req, _ := http.NewRequest("GET", "/api/v1/device/detail/"+deviceId, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-token",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InI1NjYzMzU5LTU0ZjMtZTIzYi04MDFiLWUzYjgyNDQ2Mzc4ZSIsImVtYWlsIjoic3VwZXJAc3VwZXIuY24iLCJjcmVhdGVfdGltZSI6IjIwMjQtMDMtMDZUMTQ6NDE6MTEuNzE1MDYwNCswODowMCIsImF1dGhvcml0eSI6IlNZU19BRE1JTiIsInRlbmFudF9pZCI6ImFhYWFhYWFhIn0.q2WI_eQ0837jAqCkE-Tj27IZ5C7qpYCaP9lJm1qVF2k",
	)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)

	assert.Equal(t, http.StatusOK, w.Code)
	t.Log("Get device by ID response:" + string(responseData))

	// Parse response to verify device details
	var response map[string]interface{}
	err := json.Unmarshal(responseData, &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])

	// Verify device data
	data := response["data"].(map[string]interface{})
	assert.Equal(t, deviceId, data["id"])
	assert.Equal(t, "test-get-device", data["name"])
	assert.Equal(t, "A", data["access_way"])
}
