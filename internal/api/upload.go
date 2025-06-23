package api

import (
	"crypto/md5"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Thingsly/backend/pkg/common"
	"github.com/Thingsly/backend/pkg/errcode"
	"github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UpLoadApi struct{}

// Define File Upload Configuration
const (
	BaseUploadDir = "./files/"
	OtaPath       = "./api/v1/ota/download/files/"
	MaxFileSize   = 500 << 20 // 500MB
)

// UpFile handles file upload
// @Summary File Upload
// @Description File Upload
// @Tags     File Upload
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} map[string]interface{}
// @Router   /api/v1/file/up [post]
func (*UpLoadApi) UpFile(c *gin.Context) {
	// Check if the file is empty
	file, err := c.FormFile("file")
	if err != nil || file == nil {
		c.Error(errcode.New(errcode.CodeFileEmpty))
		return
	}

	// Validate the file type
	fileType := c.PostForm("type")
	if fileType == "" {
		c.Error(errcode.New(errcode.CodeFileEmpty))
		return
	}

	// Validate the file size
	if file.Size > MaxFileSize {
		c.Error(errcode.WithVars(errcode.CodeFileTooLarge, map[string]interface{}{
			"max_size":     "500MB",
			"current_size": fmt.Sprintf("%.2fMB", float64(file.Size)/(1<<20)),
		}))
		return
	}

	// Sanitize the file name
	filename := sanitizeFilename(file.Filename)

	// Validate the file type
	if err := validateFileType(filename, fileType); err != nil {
		c.Error(errcode.WithVars(errcode.CodeFileTypeMismatch, map[string]interface{}{
			"expected_type": fileType,
			"actual_type":   filepath.Ext(filename),
		}))
		return
	}

	// Generate the file path
	uploadDir, fileName, err := generateFilePath(fileType, file.Filename)
	if err != nil {
		c.Error(errcode.WithVars(errcode.CodeFilePathGenError, map[string]interface{}{
			"error":     err.Error(),
			"file_type": fileType,
			"filename":  file.Filename,
		}))
		return
	}

	// Save the file
	filePath, err := saveFile(c, file, uploadDir, fileName, fileType)
	if err != nil {
		c.Error(errcode.WithVars(errcode.CodeFileSaveError, map[string]interface{}{
			"error":      err.Error(),
			"upload_dir": uploadDir,
			"filename":   fileName,
		}))
		return
	}

	c.Set("data", map[string]interface{}{
		"path": filePath,
	})
}

// generateFilePath generates a secure file path. The path format is: ./files/{type}/{2025-04-10}/
func generateFilePath(fileType, filename string) (string, string, error) {
	// 1. Validate if fileType contains illegal characters
	if strings.ContainsAny(fileType, "./\\") {
		return "", "", errcode.New(errcode.CodeFilePathGenError)
	}

	// 2. Generate the date directory
	dateDir := time.Now().Format("2006-01-02")

	// 3. Use filepath.Clean to clean and validate the path
	uploadDir := filepath.Clean(filepath.Join(BaseUploadDir, fileType, dateDir))
	absUploadDir, err := filepath.Abs(uploadDir)
	if err != nil {
		return "", "", errcode.WithVars(errcode.CodeFilePathGenError, map[string]interface{}{
			"error": "invalid path",
		})
	}

	// Get the absolute path of the base directory
	absBaseDir, err := filepath.Abs(BaseUploadDir)
	if err != nil {
		return "", "", errcode.WithVars(errcode.CodeFilePathGenError, map[string]interface{}{
			"error": "invalid base path",
		})
	}

	// Ensure the generated path is within the base directory
	if !strings.HasPrefix(absUploadDir, absBaseDir) {
		return "", "", errcode.WithVars(errcode.CodeFilePathGenError, map[string]interface{}{
			"error": "path traversal detected",
		})
	}

	// 4. Create the directory
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", "", errcode.WithVars(errcode.CodeFilePathGenError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	// 5. Generate a unique filename
	randomStr, err := common.GenerateRandomString(16)
	if err != nil {
		return "", "", errcode.WithVars(errcode.CodeFilePathGenError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	timeStr := time.Now().Format("20060102150405")
	hashStr := fmt.Sprintf("%x", md5.Sum([]byte(timeStr+randomStr)))
	fileName := hashStr + strings.ToLower(filepath.Ext(filename))

	return uploadDir, fileName, nil
}

// saveFile saves the file and returns the path
func saveFile(c *gin.Context, file *multipart.FileHeader, uploadDir, fileName, fileType string) (string, error) {
	fullPath := filepath.Join(uploadDir, fileName)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	// Special handling for upgrade package path
	if fileType == "upgradePackage" {
		return "./" + filepath.Join(OtaPath, fileType, time.Now().Format("2006-01-02"), fileName), nil
	}

	return "./" + fullPath, nil
}

// sanitizeFilename sanitizes the filename
func sanitizeFilename(filename string) string {
	ext := filepath.Ext(filename)
	nameOnly := strings.TrimSuffix(filepath.Base(filename), ext)

	reg := regexp.MustCompile(`[^a-zA-Z0-9-_]+`)
	sanitized := reg.ReplaceAllString(nameOnly, "_")

	if sanitized == "" || sanitized == "_" {
		sanitized = fmt.Sprintf("file_%d", time.Now().Unix())
	}

	return sanitized + strings.ToLower(ext)
}

// validateFileType validates the file type
func validateFileType(filename, fileType string) error {
	if err := utils.CheckPath(fileType); err != nil {
		return fmt.Errorf("invalid file path: %v", err)
	}

	if !utils.ValidateFileType(filename, fileType) {
		return errors.New("file type is not allowed")
	}

	return nil
}
