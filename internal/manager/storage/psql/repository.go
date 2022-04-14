package psql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rasulov-emirlan/netping-manager/internal/manager"
)

type repository struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) (manager.Repository, error) {
	if conn == nil {
		return nil, errors.New("repository: connection cannot be nil")
	}
	return &repository{
		conn: conn,
	}, nil
}

const sqlCreateSocket = ``

func (r *repository) CreateSocket(ctx context.Context, s manager.Socket) (*manager.Socket, error) {
	panic("not implemented")
}

func (r *repository) UpdateSocket(ctx context.Context, s manager.Socket) (*manager.Socket, error) {
	panic("not implemented")
}

func (r *repository) DeleteSocket(ctx context.Context, socketID int) error  {
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
