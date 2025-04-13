package initialize

import (
	"fmt"
	"log"
	"path/filepath"

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
		levelColor, levelText = "34", "DEBUG"
	case logrus.InfoLevel:
		levelColor, levelText = "36", "INFO "
	case logrus.WarnLevel:
		levelColor, levelText = "33", "WARN "
	case logrus.ErrorLevel:
		levelColor, levelText = "31", "ERROR"
	case logrus.FatalLevel, logrus.PanicLevel:
		levelColor, levelText = "31", "FATAL"
	default:
		levelColor, levelText = "0", "UNKNOWN"
	}

	var fileAndLine string
	if entry.HasCaller() {
		dir := filepath.Dir(entry.Caller.File)
		fileAndLine = fmt.Sprintf("%s/%s:%d", filepath.Base(dir), filepath.Base(entry.Caller.File), entry.Caller.Line)
	}

	msg := fmt.Sprintf("\033[1;%sm%s\033[0m \033[4;1;%sm[%s]\033[0m \033[1;%sm[%s]\033[0m %s\n",
		levelColor, levelText,
		levelColor, entry.Time.Format("2006-01-02 15:04:05.9999"),
		levelColor, fileAndLine,
		entry.Message,
	)

	return []byte(msg), nil
}
func LogInIt() {

	logrus.SetReportCaller(true)
	logrus.SetFormatter(&customFormatter{logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	}})

	var logLevels = map[string]logrus.Level{
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

	log.Println("Logrus setup complete")
}
