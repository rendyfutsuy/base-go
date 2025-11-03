package utils

import (
	"log"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var (
	ConfigVars *koanf.Koanf
)

func InitConfig(path string) {

	ConfigVars = koanf.New(".")
	if err := ConfigVars.Load(file.Provider(path), json.Parser()); err != nil {
		// Handle error loading configuration
		// Logger.Error("Koanf : " + err.Error())
		if err := ConfigVars.Load(file.Provider("config.json"), json.Parser()); err != nil {
			log.Fatalf("error loading config: %v", err)
			// Logger.Error("Koanf : " + err.Error())
		}
	}
}

func GetString(key string) string {
	return ConfigVars.String(key)
}

func GetInt(key string) int {
	return ConfigVars.Int(key)
}

func GetBool(key string) bool {
	return ConfigVars.Bool(key)
}
