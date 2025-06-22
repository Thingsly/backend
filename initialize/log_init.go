package initialize

import (
	"fmt"
	"log"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type customFormatter struct {
	logrus.TextFormatter
}

func (*customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor string
	var levelText string
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor, levelText = "34", "DEBUG" // Màu xanh dương
	case logrus.InfoLevel:
		levelColor, levelText = "36", "INFO " // Màu tím
	case logrus.WarnLevel:
		levelColor, levelText = "33", "WARN " // Màu vàng
	case logrus.ErrorLevel:
		levelColor, levelText = "31", "ERROR" // Màu đỏ
	case logrus.FatalLevel, logrus.PanicLevel:
		levelColor, levelText = "31", "FATAL" // Màu đỏ
	default:
		levelColor, levelText = "0", "UNKNOWN" // Không màu
	}

	var fileAndLine string
	if entry.HasCaller() {
		filePath := entry.Caller.File
		// Chỉ lấy phần đường dẫn từ thư mục backend trở đi
		// Input: /Users/hantdev/soict/backend/internal/service/user.go
		// Output: /internal/service/user.go:123
		if idx := strings.Index(filePath, "backend"); idx != -1 {
			filePath = filePath[idx+len("backend"):]
			if strings.HasPrefix(filePath, "/") || strings.HasPrefix(filePath, "\\") {
				filePath = "." + filePath
			} else {
				filePath = "./" + filePath
			}
		}
		fileAndLine = fmt.Sprintf("%s:%d", filePath, entry.Caller.Line)
	}

	// Output: INFO  [2024-01-15 14:30:25.1234] User login successful [./internal/service/auth.go:78]
	msg := fmt.Sprintf("\033[1;%sm%s\033[0m \033[4;1;%sm[%s]\033[0m %s \033[1;%sm[%s]\033[0m\n",
		levelColor, levelText,
		levelColor, entry.Time.Format("2006-01-02 15:04:05.9999"),
		entry.Message,
		levelColor, fileAndLine,
	)

	return []byte(msg), nil
}

func LogInIt() error {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&customFormatter{logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	}})

	logLevels := map[string]logrus.Level{
		"panic": logrus.PanicLevel,
		"fatal": logrus.FatalLevel,
		"error": logrus.ErrorLevel,
		"warn":  logrus.WarnLevel,
		"info":  logrus.InfoLevel,
		"debug": logrus.DebugLevel,
		"trace": logrus.TraceLevel,
	}

	levelStr := viper.GetString("log.level")
	if level, ok := logLevels[levelStr]; ok {
		logrus.SetLevel(level)
	} else {
		logrus.Error("Invalid log level in config, setting to default level")
		logrus.SetLevel(logrus.InfoLevel)
	}

	log.Println("Logrus setup completed")
	return nil
}
