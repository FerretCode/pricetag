package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Permissions []string

func (p Permissions) Include(name string) bool {
	for i := range p {
		if name == p[i] {
			return true
		}
	}
	return false
}

type PermissionModel struct {
	db *sql.DB
}

func (m PermissionModel) AllCodes() (Permissions, error) {
	query := `
		SELECT code
		FROM permissions;`

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions Permissions
	for rows.Next() {
		var p string
		err := rows.Scan(&p)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (m PermissionModel) GetAllForUser(userID int) (Permissions, error) {
	query := `
		SELECT permissions.code
		FROM permissions
		INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
		INNER JOIN users ON users_permissions.user_id = users.id
		WHERE users.id = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions Permissions
	for rows.Next() {
		var p string
		err := rows.Scan(&p)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (m PermissionModel) SetForUser(userID int, codes ...string) error {
	placeholders := make([]string, len(codes))
	for i := range codes {
		placeholders[i] = "?"
	}

	placeholdersStr := strings.Join(placeholders, ",")

	query := fmt.Sprintf(`
		INSERT INTO users_permissions (user_id, permission_id)
		SELECT ?, permissions.id FROM permissions
		WHERE permissions.code IN (%s)
		ON CONFLICT (user_id, permission_id) DO NOTHING;

		DELETE FROM users_permissions
		WHERE user_id = ?
		AND permission_id NOT IN (
			SELECT id FROM permissions
			WHERE permissions.code IN (%s)
		);`, placeholdersStr, placeholdersStr)

	args := []any{userID}
	for _, code := range codes {
		args = append(args, code)
	}

	// Again for the deletion
	args = append(args, userID)
	for _, code := range codes {
		args = append(args, code)
	}

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, args...)

	return err
}
