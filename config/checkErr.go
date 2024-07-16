package config

import (
	"github.com/mirjalilova/authService/config/logger"
)

type ErrorManager struct {
	logger *logger.Logger
}

func NewErrorManager(logger *logger.Logger) *ErrorManager {
	return &ErrorManager{logger: logger}
}

func (e *ErrorManager) CheckErr(err error) {
	if err != nil {
		e.logger.ERROR.Panic(err.Error())
	}
}
