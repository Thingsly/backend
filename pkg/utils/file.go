package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// User input path security check
func CheckPath(param string) error {
	if count := strings.Count(param, "."); count > 0 {
		return errors.New("path cannot contain illegal character \".\"")
	}
	if count := strings.Count(param, "/"); count > 0 {
		return errors.New("path cannot contain illegal character \"/\"")
	}
	if count := strings.Count(param, "\\"); count > 0 {
		return errors.New("path cannot contain illegal character \"\\\"")
	}
	return nil
}

// User input filename security check
func CheckFilename(param string) error {
	if count := strings.Count(param, "."); count > 1 {
		return errors.New("filename cannot contain more than one \".\"")
	}
	if count := strings.Count(param, "/"); count > 0 {
		return errors.New("filename cannot contain illegal character \"/\"")
	}
	if count := strings.Count(param, "\\"); count > 0 {
		return errors.New("filename cannot contain illegal character \"\\\"")
	}
	return nil
}

func FileSign(filePath string, sign string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	if sign == "MD5" {
		hash := md5.New()
		_, _ = io.Copy(hash, file)
		return hex.EncodeToString(hash.Sum(nil)), nil
	} else {
		hash := sha256.New()
		_, _ = io.Copy(hash, file)
		return hex.EncodeToString(hash.Sum(nil)), nil
	}

}

var (
	allowExtMap = map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".svg": true,
		".ico": true, ".gif": true, ".xlsx": true, ".xls": true, ".csv": true,
	}

	allowUpgradePackageMap = map[string]bool{
		".bin": true, ".tar": true, ".gz": true, ".zip": true,
		".gzip": true, ".apk": true, ".dav": true, ".pack": true,
	}

	allowImportBatchMap = map[string]bool{
		".xlsx": true, ".xls": true, ".csv": true,
	}
)

func ValidateFileType(filename, fileType string) bool {
	ext := strings.ToLower(path.Ext(filename))

	switch fileType {
	case "upgradePackage":
		return allowUpgradePackageMap[ext]
	case "importBatch":
		return allowImportBatchMap[ext]
	case "d_plugin":

		return true
	default:
		return allowExtMap[ext]
	}
}

func ValidateFileExtension(filename string, allowedExts []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return false
	}

	ext = strings.TrimPrefix(ext, ".")

	for _, allowedExt := range allowedExts {
		if strings.ToLower(allowedExt) == ext {
			return true
		}
	}

	return false
}

func SanitizeFilename(filename string) string {

	filename = filepath.Base(filename)

	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	reg := regexp.MustCompile(`[^\w\-\.]`)
	nameWithoutExt = reg.ReplaceAllString(nameWithoutExt, "_")

	nameWithoutExt = handleSpecialFilenames(nameWithoutExt)

	if len(nameWithoutExt) > 200 {
		nameWithoutExt = nameWithoutExt[:200]
	}

	if strings.HasPrefix(nameWithoutExt, ".") {
		nameWithoutExt = "_" + nameWithoutExt
	}

	ext = reg.ReplaceAllString(ext, "_")

	sanitizedName := nameWithoutExt + strings.ToLower(ext)

	if sanitizedName == "" {
		return "unnamed_file"
	}

	return sanitizedName
}

func handleSpecialFilenames(filename string) string {

	lowerName := strings.ToLower(filename)

	specialNames := map[string]bool{
		"con": true, "prn": true, "aux": true, "nul": true,
		"com1": true, "com2": true, "com3": true, "com4": true,
		"lpt1": true, "lpt2": true, "lpt3": true, "lpt4": true,
	}

	if specialNames[lowerName] {
		return "_" + filename
	}

	return filename
}
