package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	model "github.com/HustIoTPlatform/backend/internal/model"
	query "github.com/HustIoTPlatform/backend/internal/query"
	utils "github.com/HustIoTPlatform/backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
)

var allowedFileExts = []string{
	"jpg", "jpeg", "png", "pdf", "doc", "docx", "xlsx", "xls", "zip", "rar", "tar", "gz", "7z",
}

// OperationLogs Middleware for logging operations
func OperationLogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only log for methods that modify resources (e.g., POST, PUT, DELETE)
		if !isModifyMethod(c.Request.Method) {
			c.Next()
			return
		}

		// Log the request method and URL path
		logrus.Info("Started processing request:", c.Request.URL.Path, "Method:", c.Request.Method)
		requestMessage, _ := processRequestBody(c) // Process and log the request body
		logrus.Info("Request body:", requestMessage)

		// Create a custom response body writer to capture response body
		writer := newResponseBodyWriter(c)
		c.Writer = writer

		// Record the start time
		start := time.Now().UTC()
		c.Next()

		// Measure the processing time
		cost := time.Since(start).Milliseconds()

		// Log the status code, processing time, and response body
		logrus.Info("Request processing completed, Status Code:", c.Writer.Status(), "Duration (ms):", cost)
		logrus.Info("Response body size:", writer.body.Len())
		logrus.Info("Response body:", writer.body.String())

		// Save the operation log
		saveOperationLog(c, start, cost, requestMessage, writer.body.String())
	}
}


func isModifyMethod(method string) bool {
	return method == http.MethodPost ||
		method == http.MethodPut ||
		method == http.MethodDelete
}

func processRequestBody(c *gin.Context) (string, string) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error("Failed to read the request body:", err)
		return "", ""
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	requestMessage := string(body)
	if strings.Contains(c.Request.URL.Path, "file/up") {
		requestMessage = handleFileUpload(c)
	}

	return requestMessage, requestMessage
}

func handleFileUpload(c *gin.Context) string {
    // Retrieve the uploaded file from the request form
    file, err := c.FormFile("file")
    if err != nil {
        return "" // If there's an error, return an empty string
    }

    // Retrieve file type from the form, default to "unknown" if not specified
    fileType := c.PostForm("type")
    if fileType == "" {
        fileType = "unknown"
    }

    // 1. Extract a safe file name, removing path
    baseFileName := filepath.Base(file.Filename)

    // 2. Sanitize the file name (remove dangerous characters, etc.)
    filename := utils.SanitizeFilename(baseFileName)

    // 3. Second check on filename's safety (ensuring it's a valid file path)
    if !filepath.IsLocal(filename) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Illegal file name"}) // Respond with an error if filename is not safe
        return ""
    }

    // 4. Validate the file type/extension
    if !utils.ValidateFileExtension(filename, allowedFileExts) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Disallowed file type"}) // Respond with an error if the file type is not allowed
        return ""
    }

    // Return a formatted string with the file type and filename
    return fmt.Sprintf("%s:%s", fileType, filename)
}

func saveOperationLog(c *gin.Context, start time.Time, cost int64, requestMsg, responseMsg string) {
    // Check if claims exist in the context (this holds user info)
    claims, exists := c.Get("claims")
    if !exists {
        logrus.Info("No user info found, skipping operation log")
        return
    }

    // Type assertion to get user claims from the context
    userClaims, ok := claims.(*utils.UserClaims)
    if !ok {
        logrus.Info("Incorrect user info type, skipping operation log")
        return
    }

    // Check if tenantID is empty
    if userClaims.TenantID == "" {
        logrus.Info("TenantID is empty, skipping operation log")
        return
    }

    // Prepare the operation log
    path := c.Request.URL.Path
    log := &model.OperationLog{
        ID:              uuid.New(),
        IP:              c.ClientIP(),
        Path:            &path,
        UserID:          userClaims.ID,
        Name:            &c.Request.Method,
        CreatedAt:       start,
        Latency:         &cost,
        RequestMessage:  &requestMsg,
        ResponseMessage: &responseMsg,
        TenantID:        userClaims.TenantID,
    }

    // Save the log to the database
    query.OperationLog.Create(log)
}


type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func newResponseBodyWriter(c *gin.Context) responseBodyWriter {
	return responseBodyWriter{
		ResponseWriter: c.Writer,
		body:           &bytes.Buffer{},
	}
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
