package configtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/test/library-app/internal/config"
)

func TestLoadConfig(t *testing.T) {
	err := config.LoadConfig()
	assert.Nil(t, err)
}
