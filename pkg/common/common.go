package common

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	constant "github.com/Thingsly/backend/pkg/constant"

	"github.com/pkg/errors"
)

func CheckEmpty(str string) bool {
	return str == constant.EMPTY
}

func GetMessageID() string {

	// timestamp := time.Now().UnixNano() // Sử dụng nano seconds
	// random := rand.Intn(1000)          // Thêm random
	// return fmt.Sprintf("%d%03d", timestamp%10000000, random)

	timestamp := time.Now().Unix()

	// Convert timestamp to string base 10
	timestampStr := strconv.FormatInt(timestamp, 10)

	// Get the last 7 digits of the timestamp string
	messageID := timestampStr[len(timestampStr)-7:]

	return messageID
}

// JsonToString
func JsonToString(any any) (string, error) {
	data, err := json.Marshal(any)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func GetErrors(err error, message string) error {
	return errors.WithMessage(err, message)
}

// GetResponsePayload return response payload for API response
func GetResponsePayload(method string, err error) []byte {
	if err != nil {
		data := map[string]interface{}{
			"result":  1,
			"errcode": "000",
			"message": err.Error(),
			"ts":      time.Now().Unix(),
		}
		res, _ := json.Marshal(data)
		return res
	}
	data := map[string]interface{}{
		"result":  0,
		"message": "success",
		"ts":      time.Now().Unix(),
	}
	if method != "" {
		data["method"] = method
	}
	res, _ := json.Marshal(data)
	return res
}

func StringSpt(str string) *string {
	return &str
}

func IsStringEmpty(str *string) bool {
	return str == nil || *str == ""
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b), nil
}

var ErrNoRows = errors.New("record not found")

func GetRandomNineDigits() (string, error) {

	min := big.NewInt(100000000)
	max := big.NewInt(999999999)

	diff := new(big.Int).Sub(max, min)
	diff = diff.Add(diff, big.NewInt(1))

	n, err := rand.Int(rand.Reader, diff)
	if err != nil {
		return "", fmt.Errorf("failed to generate random number: %v", err)
	}

	n = n.Add(n, min)

	return n.String(), nil
}

func GenerateNumericCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}

	code := make([]byte, length)

	for i := 0; i < length; i++ {

		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %v", err)
		}

		code[i] = byte(num.Int64() + '0')
	}

	return string(code), nil
}
