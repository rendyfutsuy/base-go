package utils

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var (
	ConfigVars *koanf.Koanf
)

func InitConfig(path string) {
	ConfigVars = koanf.New(".")

	loadEnvFile(path)
	loadEnvFile(".env")

	if err := ConfigVars.Load(file.Provider("config.json"), json.Parser()); err != nil {
		log.Printf("config.json not loaded: %v", err)
	}

	applyEnvToKoanf()

	InitEmailDomainRegex()
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

func loadEnvFile(path string) {
	if path == "" {
		return
	}
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	re := regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_\.]*)\s*=\s*(.*)\s*$`)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "#") || strings.TrimSpace(line) == "" {
			continue
		}
		m := re.FindStringSubmatch(line)
		if len(m) == 3 {
			key := strings.TrimSpace(m[1])
			val := strings.TrimSpace(m[2])
			if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") && len(val) >= 2 {
				val = strings.TrimSuffix(strings.TrimPrefix(val, "\""), "\"")
			}
			if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") && len(val) >= 2 {
				val = strings.TrimSuffix(strings.TrimPrefix(val, "'"), "'")
			}
			os.Setenv(key, val)
		}
	}
}

func applyEnvToKoanf() {
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) != 2 {
			continue
		}
		name := parts[0]
		value := parts[1]
		key := envToKoanfKey(name)
		if key == "" {
			continue
		}
		ConfigVars.Set(key, parseEnvValue(value))
	}
}

func envToKoanfKey(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, "__", ".")
	return s
}

func parseEnvValue(v string) interface{} {
	lower := strings.ToLower(v)
	if lower == "true" {
		return true
	}
	if lower == "false" {
		return false
	}
	if iRe.MatchString(v) {
		return toInt(v)
	}
	if fRe.MatchString(v) {
		return toFloat(v)
	}
	return v
}

var iRe = regexp.MustCompile(`^[-+]?\d+$`)
var fRe = regexp.MustCompile(`^[-+]?\d*\.\d+$`)

func toInt(s string) int64 {
	sign := 1
	if strings.HasPrefix(s, "-") {
		sign = -1
		s = s[1:]
	}
	var res int64
	for i := 0; i < len(s); i++ {
		c := s[i] - '0'
		res = res*10 + int64(c)
	}
	return int64(sign) * res
}

func toFloat(s string) float64 {
	sign := 1.0
	if strings.HasPrefix(s, "-") {
		sign = -1.0
		s = s[1:]
	}
	if idx := strings.IndexByte(s, '.'); idx >= 0 {
		intPart := s[:idx]
		fracPart := s[idx+1:]
		var iVal int64
		for i := 0; i < len(intPart); i++ {
			c := intPart[i] - '0'
			iVal = iVal*10 + int64(c)
		}
		var fVal float64
		pow := 1.0
		for i := 0; i < len(fracPart); i++ {
			c := float64(fracPart[i] - '0')
			fVal = fVal*10 + c
			pow *= 10
		}
		return sign * (float64(iVal) + fVal/pow)
	}
	return sign * float64(toInt(s))
}
