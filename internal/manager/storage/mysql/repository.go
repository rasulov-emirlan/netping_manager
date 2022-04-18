package mysql

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rasulov-emirlan/netping-manager/internal/manager"
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

const createSocketSQL = `
	INSERT INTO sockets (
		name, mib_address, netping_address, socket_type_id
	)
	VALUES( ?, ?, ?, ? );
`

func (r *repository) CreateSocket(ctx context.Context, s manager.Socket) (*manager.Socket, error) {
	err := r.conn.QueryRow(createSocketSQL, s.Name, s.SNMPmib, s.SNMPaddress, s.ObjectType).Scan(&s.ID)
	return &s, err
}

const updateSocketSQL = `
	UPDATE sockets SET name = ?, mib_addre	ss = ?, netping_address = ?, socket_type_id = ?
		WHERE id = ?;
`

func (r *repository) UpdateSocket(ctx context.Context, s manager.Socket) (*manager.Socket, error) {
	return &s, r.conn.QueryRow(updateSocketSQL, s.Name, s.SNMPmib, s.SNMPaddress, s.ObjectType, s.ID).Err()
}

const deleteSocketSQL = `
	DELETE FROM sockets WHERE id = ?;
`

func (r *repository) DeleteSocket(ctx context.Context, socketID int) error {
	return r.conn.QueryRow(deleteSocketSQL, socketID).Err()
}

const findSocketsByLocationSQL = `
	SELECT id, name, mib_address, socket_type_id FROM sockets WHERE netping_address = ?;
`

func (r *repository) FindSocketsByLocation(ctx context.Context, locationAddress string) ([]*manager.Socket, error) {
	rows, err := r.conn.Query(findSocketByIDsql, locationAddress)
	if err != nil {
		return nil, err
	}

	var (
		sockets []*manager.Socket
		id, socketType int
		name, mibAddress string
	)

	for rows.Next() {
		if err := rows.Scan(
			&id, &name, &mibAddress, &socketType,
		); err != nil {
			return nil, err
		}
		sockets = append(sockets, &manager.Socket{
			ID: id,
			Name: name,
			SNMPaddress: locationAddress,
			SNMPmib: mibAddress,
			ObjectType: socketType,
		})
	}
	return sockets, nil
}

const findSocketByIDsql = `
	SELECT name, mib_address, netping_address, socket_type_id FROM sockets WHERE id = ?;
`

func (r *repository) FindSocketByID(ctx context.Context, socketID int) (*manager.Socket, error) {
	socket := &manager.Socket{
		ID: socketID,
	}
	err := r.conn.QueryRow(findSocketByIDsql, socketID).Scan(&socket.Name, &socket.SNMPmib, &socket.SNMPaddress, &socket.ObjectType)
	return socket, err
}

const listAllSocketsSQL = `
	SELECT netping_address, id, name, mib_address, socket_type_id from sockets group by netping_address ORDER BY netping_address ASC;
`

func (r *repository) ListAllSockets(ctx context.Context) ([]*manager.Location, error) {
	rows, err := r.conn.Query(listAllSocketsSQL)
	if err != nil {
		return nil, err
	}

	var (
		locationsMap = make(map[string][]*manager.Socket)
		netpingAddress, name, mibAddress string
		id, socketType int
	)

	for rows.Next() {
		if err := rows.Scan(
			&netpingAddress,
			&id,
			&name,
			&mibAddress,
			&socketType,
		); err != nil {
			return nil, err
		}
		locationsMap[netpingAddress] = append(locationsMap[netpingAddress], &manager.Socket{
			ID: id,
			Name: name,
			SNMPaddress: netpingAddress,
			SNMPmib: mibAddress,
		})
	}

	locations := make([]*manager.Location, len(locationsMap))
	index := 0
	for i, v := range locationsMap {
		locations[index] = &manager.Location{
			SNMPaddress: i,
			Sockets: v,
		}
		index++
	}

	return locations, nil
}