package types

import "time"

type User struct {
	ID           int    `db:"ID"`
	Username     string `db:"Username"`
	Password     string `db:"Password"`
	PermissionID int    `db:"PermissionID"`
}

type Permission struct {
	ID     int `db:"ID"`
	UserID int `db:"UserID"`

	Admin            bool `db:"Admin"`
	ManageServices   bool `db:"ManageServices"`
	ManageTags       bool `db:"ManageTags"`
	ManageForwarding bool `db:"ManageForwarding"`
	ViewLogs         bool `db:"ViewLogs"`
}

type Service struct {
}

type Tag struct {
}

type Log struct {
	Message   string    `json:"message"`
	Level     string    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
}

type Error struct {
	Status int
	Error  string
}
