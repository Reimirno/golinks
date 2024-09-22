package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name  string
		debug bool
	}{
		{"prod", false},
		{"dev", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Initialize(tt.debug)
			assert.NoError(t, err)
		})
	}
}

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name       string
		debug      bool
		loggerName string
	}{
		{"prod", false, "l1"},
		{"dev", true, "l2"},
		{"dev", true, "l3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerSet = nil
			err := Initialize(tt.debug)
			assert.NoError(t, err)

			logger := NewLogger(tt.loggerName)
			assert.NotNil(t, logger)
			assert.Equal(t, tt.loggerName, logger.Desugar().Name())
		})
	}
}
