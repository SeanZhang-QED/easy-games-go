package models

import (
	"time"
)

type Session struct {
	Email  string 
	LastActivity time.Time 
}
