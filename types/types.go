package types

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

type Error struct {
	Status int
	Error  string
}
