package models

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var Cache = cache.New(1*time.Hour, 10*time.Minute)