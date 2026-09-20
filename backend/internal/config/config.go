package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Addr        string
	DatabaseURL string
	JWTSecret   string
	JWTAccessTTL string
	StorageDir  string
	AdminEmail  string
	SeedPassword string
	RedisAddr   string
	TelegramBotToken string
	TelegramChatIDs  []int64
}

func Load() Config {
	return Config{
		Addr:            getenv("ADDR", ":8080"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		JWTAccessTTL:    getenv("JWT_ACCESS_TTL", "72h"),
		StorageDir:      getenv("STORAGE_DIR", "storage/proofs"),
		AdminEmail:      getenv("ADMIN_EMAIL", "admin@example.com"),
		SeedPassword:    getenv("SEED_ADMIN_PASSWORD", "admin12345"),
		RedisAddr:       getenv("REDIS_ADDR", "localhost:6379"),
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatIDs: splitInt64(os.Getenv("TELEGRAM_CHAT_IDS")),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func splitInt64(s string) []int64 {
	var out []int64
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if n, err := strconv.ParseInt(p, 10, 64); err == nil {
			out = append(out, n)
		}
	}
	return out
}
