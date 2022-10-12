package internal

import "gorm.io/gorm"

type RequestInfo struct {
	gorm.Model
	RemoteAddress string
	UserAgent     string
}
