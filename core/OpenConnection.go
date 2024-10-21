package core

import (
	"sync"

	"gorm.io/gorm"
)

type OpenConnection struct {
	connection *gorm.DB
	mux        sync.Mutex
}
