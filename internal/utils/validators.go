package utils

import (
	"go.uber.org/zap"
	"net/mail"
	"time"
)

func IsEmailValid(email string) (result0 bool) {
	__logParams := map[string]any{"email": email}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "IsEmailValid"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "IsEmailValid"), zap.Any("params", __logParams))
	_, err := mail.ParseAddress(email)
	result0 = err == nil
	return
}
