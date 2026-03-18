package services

import (
	"encoding/json"
	"fmt"
	"time"

	"go-rest/config"
	"go-rest/utils"
)

// SetSession simpan session user ke Redis
func SetSession(user utils.CustomerSession, duration time.Duration) error {
	key := fmt.Sprintf("session:%d", user.UserID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return config.RDB.Set(config.Ctx, key, data, duration).Err()
}

// GetSession ambil session dari Redis
func GetSession(userID int64) (*utils.CustomerSession, error) {
	key := fmt.Sprintf("session:%d", userID)
	val, err := config.RDB.Get(config.Ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var session utils.CustomerSession
	if err := json.Unmarshal([]byte(val), &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// DeleteSession hapus session dari Redis
func DeleteSession(userID int64) error {
	key := fmt.Sprintf("session:%d", userID)
	return config.RDB.Del(config.Ctx, key).Err()
}
