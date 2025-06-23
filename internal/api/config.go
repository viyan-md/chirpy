package api

import (
	"sync/atomic"

	"github.com/viyan-md/chirpy/internal/database"
)

type APIConfig struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
	Platform       string
}
