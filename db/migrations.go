package db

import (
	"database/sql"
)

func RunMigrations(db *sql.DB) error {
	createUserTableQuery := `
	CREATE TABLE User (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Username TEXT NOT NULL UNIQUE,
		Password TEXT NOT NULL,
		PermissionID INTEGER,
		FOREIGN KEY (PermissionID) REFERENCES Permission(ID)
	);`

	createPermissionsQuery := `
	CREATE TABLE Permission (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		UserID INTEGER NOT NULL,
		Admin BOOLEAN NOT NULL DEFAULT 0,
		ManageServices BOOLEAN NOT NULL DEFAULT 0,
		ManageTags BOOLEAN NOT NULL DEFAULT 0,
		ManageForwarding BOOLEAN NOT NULL DEFAULT 0,
		ViewLogs BOOLEAN NOT NULL DEFAULT 0,
		FOREIGN KEY (UserID) REFERENCES User(ID)
	);
	`

	var errors []error

	_, err := db.Query(createUserTableQuery)
	errors = append(errors, err)

	_, err = db.Query(createPermissionsQuery)
	errors = append(errors, err)

	for _, err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}
