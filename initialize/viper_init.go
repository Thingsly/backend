package initialize

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

func ViperInit(path string) error {
	viper.SetEnvPrefix("THINGSLY")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if path != "" {
		viper.SetConfigFile(path)
	} else {
		viper.SetConfigName("./configs/conf")
	}

	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %s", err)
	}
	log.Println("Viper has finished loading the conf.yml configuration file...")
	return nil
}
