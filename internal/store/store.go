package store

import (
	"fmt"

	"github.com/test/library-app/internal/config"
	"github.com/test/library-app/internal/constants"
	"github.com/test/library-app/internal/store/local"
)

type Store interface {
}

func NewStore() (Store, error) {
	switch config.CommonConfig.StoreType {
	case constants.LocalStore:
		return local.InitLocalStore()
	default:
		return nil, fmt.Errorf("unknown Store configured: %v", config.CommonConfig.StoreType)
	}
}
