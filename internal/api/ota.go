package api

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/errcode"
	"github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/howeyc/crc16"
	"github.com/sirupsen/logrus"
)

type OTAApi struct{}

// CreateOTAUpgradePackage
// @Router   /api/v1/ota/package [post]
func (*OTAApi) CreateOTAUpgradePackage(c *gin.Context) {
	var req model.CreateOTAUpgradePackageReq
	if !BindAndValidate(c, &req) {
		return
	}
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	err := service.GroupApp.OTA.CreateOTAUpgradePackage(&req, userClaims.TenantID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// DeleteOTAUpgradePackage
// @Router   /api/v1/ota/package/{id} [delete]
func (*OTAApi) DeleteOTAUpgradePackage(c *gin.Context) {
	id := c.Param("id")
	err := service.GroupApp.OTA.DeleteOTAUpgradePackage(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// UpdateOTAUpgradePackage
// @Router   /api/v1/ota/package/ [put]
func (*OTAApi) UpdateOTAUpgradePackage(c *gin.Context) {
	var req model.UpdateOTAUpgradePackageReq
	if !BindAndValidate(c, &req) {
		return
	}
	err := service.GroupApp.OTA.UpdateOTAUpgradePackage(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// GetOTAUpgradePackageByPage
// @Router   /api/v1/ota/package [get]
func (*OTAApi) HandleOTAUpgradePackageByPage(c *gin.Context) {
	var req model.GetOTAUpgradePackageLisyByPageReq
	if !BindAndValidate(c, &req) {
		return
	}
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	list, err := service.GroupApp.OTA.GetOTAUpgradePackageListByPage(&req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", list)
}

// CreateOTAUpgradeTask
// @Router   /api/v1/ota/task [post]
func (*OTAApi) CreateOTAUpgradeTask(c *gin.Context) {
	var req model.CreateOTAUpgradeTaskReq
	if !BindAndValidate(c, &req) {
		return
	}

	err := service.GroupApp.OTA.CreateOTAUpgradeTask(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// DeleteOTAUpgradeTask
// @Router   /api/v1/ota/task/{id} [delete]
func (*OTAApi) DeleteOTAUpgradeTask(c *gin.Context) {
	id := c.Param("id")
	err := service.GroupApp.OTA.DeleteOTAUpgradeTask(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// GetOTAUpgradeTaskByPage
// @Router   /api/v1/ota/task [get]
func (*OTAApi) HandleOTAUpgradeTaskByPage(c *gin.Context) {
	var req model.GetOTAUpgradeTaskListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}
	list, err := service.GroupApp.OTA.GetOTAUpgradeTaskListByPage(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", list)
}

// GetOTAUpgradeTaskDetailByPage
// @Router   /api/v1/ota/task/detail [get]
func (*OTAApi) HandleOTAUpgradeTaskDetailByPage(c *gin.Context) {
	var req model.GetOTAUpgradeTaskDetailReq
	if !BindAndValidate(c, &req) {
		return
	}
	list, err := service.GroupApp.OTA.GetOTAUpgradeTaskDetailListByPage(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", list)

}

// UpdateOTAUpgradeTaskStatus
// @Router   /api/v1/ota/task/detail [put]
func (*OTAApi) UpdateOTAUpgradeTaskStatus(c *gin.Context) {
	var req model.UpdateOTAUpgradeTaskStatusReq
	if !BindAndValidate(c, &req) {
		return
	}
	err := service.GroupApp.OTA.UpdateOTAUpgradeTaskStatus(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// GET /api/v1/ota/download/{filepath}
func (*OTAApi) DownloadOTAUpgradePackage(c *gin.Context) {
	filePath := "./files/upgradePackage/" + c.Param("path") + "/" + c.Param("file")

	if !utils.FileExist(filePath) {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"param_err": "file not exist",
		}))
		return
	}

	rangeHeader := c.GetHeader("Range")
	crc16Method := c.GetHeader("Crc16-Method")

	if rangeHeader == "" {
		c.File(filePath)
		return
	}

	// Send partial file content
	serveRangeFile(filePath, rangeHeader, crc16Method, c)
}

func serveRangeFile(filePath, rangeHeader, crc16Method string, c *gin.Context) {
	rangeStr := strings.Replace(rangeHeader, "bytes=", "", 1)
	rangeParts := strings.Split(rangeStr, "-")
	if len(rangeParts) != 2 {
		c.AbortWithError(http.StatusRequestedRangeNotSatisfiable, errors.New("invalid range"))
		return
	}

	start, err := strconv.ParseInt(rangeParts[0], 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusRequestedRangeNotSatisfiable, errors.New("invalid range"))
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Use named return value to ensure close errors are handled on function return
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			// Log close error
			log.Printf("Error closing file: %v", closeErr)
			// If no other error has occurred, return the close error
			if err == nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	fileSize := fileInfo.Size()

	if rangeParts[1] == "" {
		rangeParts[1] = fmt.Sprintf("%d", fileSize-1)
	}
	end, err := strconv.ParseInt(rangeParts[1], 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	if start >= fileSize || end >= fileSize {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	contentLength := end - start + 1

	c.Writer.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	c.Writer.Header().Set("Accept-Ranges", "bytes")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", contentLength))
	c.Writer.Header().Set("Content-Type", filePath[len(filePath)-3:])

	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Create a buffer to hold the file data
	buffer := make([]byte, contentLength)

	// Read the file data into the buffer
	_, err = file.Read(buffer)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var crcValue uint16
	switch crc16Method {
	case "CCITT":
		crcValue = crc16.ChecksumCCITT(buffer)
	case "MODBUS":
		crcValue = crc16.ChecksumMBus(buffer)
	default:
		crcValue = crc16.ChecksumIBM(buffer)
	}

	// Set the CRC16 value in the response header
	c.Writer.Header().Set("X-CRC16", fmt.Sprintf("%04x", crcValue))

	// Write the buffer to the response
	_, err = c.Writer.Write(buffer)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Sync the file to ensure all data is written
	if err = file.Sync(); err != nil {
		logrus.Errorf("Error syncing file: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
