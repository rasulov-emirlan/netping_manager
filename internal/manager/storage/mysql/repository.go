package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rasulov-emirlan/netping-manager/internal/manager"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type repository struct {
	conn *sql.DB
}

func NewRepository(conn *sql.DB) (manager.Repository, error) {
	if conn == nil {
		return nil, errors.New("repository: connection cannot be nil")
	}
	return &repository{
		conn: conn,
	}, nil
}

func (r *repository) CreateSocket(ctx context.Context, s manager.Socket) (*manager.Socket, error) {
	err := r.conn.QueryRow(`
	INSERT INTO sockets(mib_address, netping_address, socket_type_id)
	VALUES(?, ?, ?)
	RETURNING id;
	`, s.SNMPmib, s.SNMPaddress, s.ObjectType).Scan(&s.ID)
	return &s, err
}

func (r *repository) UpdateSocket(ctx context.Context, s manager.Socket) (*manager.Socket, error) {
	panic("not implemented")
}

func (r *repository) DeleteSocket(ctx context.Context, socketID int) error {
	panic("not implemented")
}

func (r *repository) CreateLocation(ctx context.Context, l manager.Location) (*manager.Location, error) {
	panic("not implemented")
}

func (r *repository) UpdateLocation(ctx context.Context, l manager.Location) (*manager.Location, error) {
	panic("not implemented")
}

func (r *repository) DeleteLocation(ctx context.Context, locationID int) error {
	panic("not implemented")
}

func (r *repository) ListLocations(ctx context.Context) ([]*manager.Location, error) {
	panic("not implemented")
}

