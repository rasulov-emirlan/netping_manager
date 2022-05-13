package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/rasulov-emirlan/netping-manager/internal/users"
)

type repository struct {
	conn *sql.DB
}

func NewRepository(conn *sql.DB) (*repository, error) {
	if conn == nil {
		return nil, errors.New("repository: connection cannot be nil")
	}
	return &repository{
		conn: conn,
	}, nil
}

const createSQL = `
	INSERT INTO netping_manager_users(name, password, is_admin)
		VALUES(?, ?, ?);
`

const createGetIdSQL = `
	SELECT id, created_at
	FROM netping_manager_users
		WHERE name = ?;
`

func (r *repository) Create(ctx context.Context, u users.User) (int, error) {
	if err := r.conn.QueryRow(createSQL, u.Name, u.Password, u.IsAdmin).Err(); err != nil {
		return 0, err
	}
	if err := r.conn.QueryRow(createGetIdSQL, u.Name).Scan(&u.ID, &u.CreatedAt); err != nil {
		return 0, err
	}
	return u.ID, nil
}

const readSQL = `
	SELECT name, password, is_admin
	FROM netping_manager_users
		WHERE id = ?;
`

func (r *repository) Read(ctx context.Context, userID int) (users.User, error) {
	u := users.User{
		ID: userID,
	}
	err := r.conn.QueryRow(readSQL, userID).Scan(&u.Name, &u.Password, &u.IsAdmin)
	return u, err
}

const readByNameSQL = `
	SELECT id, password, is_admin
	FROM netping_manager_users
		WHERE name = ?;
`

func (r *repository) ReadByName(ctx context.Context, name string) (users.User, error) {
	u := users.User{
		Name: name,
	}
	err := r.conn.QueryRow(readByNameSQL, name).Scan(&u.ID, &u.Password, &u.IsAdmin)
	return u, err
}

const readAllsql = `
	SELECT id, name, password, is_admin
	FROM netping_manager_users;
`

func (r *repository) ReadAll(ctx context.Context) ([]users.User, error) {
	rows, err := r.conn.Query(readAllsql)
	if err != nil {
		return nil, err
	}

	var u []users.User

	for rows.Next() {
		t := users.User{}
		if err := rows.Scan(&t.ID, &t.Name, &t.Password, &t.IsAdmin); err != nil {
			return nil, err
		}
		u = append(u, t)
	}
	return u, nil
}

const updateSQL = `
	UPDATE netping_manager_users
	SET name = ?, is_admin = ?, updated_at = ?
`

func (r *repository) Update(ctx context.Context, userID int, changeset users.User) error {
	switch changeset.Password {
	case "":
		sql := updateSQL + `
		WHERE id = ?`
		_, err := r.conn.ExecContext(ctx, sql, changeset.Name, changeset.IsAdmin, time.Now(), userID)
		return err
	default:
		sql := updateSQL + `, password = ?
		WHERE id = ?`
		_, err := r.conn.ExecContext(ctx, sql, changeset.Name, changeset.IsAdmin, time.Now(), changeset.Password, userID)
		return err
	}
}

const deleteSQL = `
	DELETE FROM netping_manager_users
		WHERE id = ?;
`

func (r *repository) Delete(ctx context.Context, userID int) error {
	return r.conn.QueryRow(deleteSQL, userID).Err()
}
