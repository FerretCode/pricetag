package types

type User struct {
	ID           int
	Username     string
	Password     string
	PermissionID int
}

type Permission struct {
	ID     int
	UserID int

	Admin            bool
	ManageServices   bool
	ManageTags       bool
	ManageForwarding bool
	ViewLogs         bool
}

type Service struct {
}

type Tag struct {
}

type Error struct {
	Status int
	Error  string
}
