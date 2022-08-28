package logger

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func InitInfoLogger(conf *viper.Viper) (*log.Logger, error) {
	logPath := strings.TrimSuffix(conf.GetString("log"), "/")
	fl, err := os.OpenFile(fmt.Sprintf("%s/%s", logPath, "app-log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
	if err != nil {
		return nil, err
	}

	return log.New(fl, "[INFO] ", log.Ldate|log.Ltime), nil
}

func InitErrorLogger(conf *viper.Viper) (*log.Logger, error) {
	logPath := strings.TrimSuffix(conf.GetString("log"), "/")
	fl, err := os.OpenFile(fmt.Sprintf("%s/%s", logPath, "app-log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
	if err != nil {
		return nil, err
	}

	return log.New(fl, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile), nil
}
