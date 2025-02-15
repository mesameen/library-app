package storetest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/test/library-app/internal/config"
	"github.com/test/library-app/internal/store"
)

func TestMain(m *testing.M) {
	config.LoadConfig()
}

func TestNewStore(t *testing.T) {
	store, err := store.NewStore()
	assert.Nil(t, err)
	assert.NotNil(t, store)
}
