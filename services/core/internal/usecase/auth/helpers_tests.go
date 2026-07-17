package auth

import (
	"core_service/internal/config"
	"time"
)

func testJWTConfig() config.JWTConfig {
	return config.JWTConfig{
		AccessSecret:  "test-access-secret",
		RefreshSecret: "test-refresh-secret",
		AccessTTL:     time.Hour,
		RefreshTTL:    time.Hour * 24,
	}
}
