package loggertest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/test/library-app/internal/config"
	"github.com/test/library-app/internal/logger"
)

func TestMain(m *testing.M) {
	config.LoadConfig()
}

func TestInitLogger(t *testing.T) {
	err := logger.InitLogger()
	assert.Nil(t, err)
}
