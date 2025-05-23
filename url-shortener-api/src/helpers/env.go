package helpers

import (
	"bufio"
	"os"
	"strings"
	"url-shortener/src/configs"
)

func LoadEnvFile(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		value = strings.Trim(value, `"'`)

		err := os.Setenv(key, value)
		if err != nil {
			panic(err)
		}
	}
}

func LoadConfigFromEnv() *configs.Config {
	return &configs.Config{
		MongoURI:                    os.Getenv("MONGO_URI"),
		MongoURLShortenerDatabase:   os.Getenv("MONGO_URL_SHORTENER_DATABASE"),
		MongoURLShortenerCollection: os.Getenv("MONGO_URL_SHORTENER_COLLECTION"),
		BaseURL:                     os.Getenv("BASE_URL"),
	}
}
