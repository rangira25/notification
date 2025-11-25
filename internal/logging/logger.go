package logging

import (
	"go.uber.org/zap"
)

func NewLogger(debug bool) *zap.Logger {
	if debug {
		l, _ := zap.NewDevelopment()
		return l
	}
	l, _ := zap.NewProduction()
	return l
}
