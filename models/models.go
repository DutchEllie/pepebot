package models

import (
	"errors"
	"time"
)


var ErrNoRecord = errors.New("models: no matching record found")

type Badword struct {
	ID       int
	Word     string
	ServerID string
	LastSaid time.Time
}

type AdminRoles struct {
	ID int
	RoleName string
	RoleID string
	GuildID string
}