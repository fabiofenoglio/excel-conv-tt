package config

import (
	"context"

	"github.com/sirupsen/logrus"
)

type WorkflowContext struct {
	Context context.Context
	Logger  *logrus.Entry
	Config  WorkflowContextConfig
}

type WorkflowContextConfig struct {
	EnableMissingOperatorsWarning bool
}
