package mysql

import (
	"context"
	"database/sql"
	"errors"
	"log"

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
	log.Println(s.Name)
	log.Println(s.SNMPmib)
	log.Println(s.SNMPaddress)
	log.Println(s.ObjectType)
	
	err := r.conn.QueryRow(createSocketSQL, s.Name, s.SNMPmib, s.SNMPaddress, s.ObjectType).Scan(&s.ID)
	return &s, err
}

func (r *repository) UpdateSocket(ctx context.Context, s manager.Socket) (*manager.Socket, error) {
	return &s, r.conn.QueryRow(`
	UPDATE sockets SET name = ?, mib_addre	ss = ?, netping_address = ?, socket_type_id = ?
	WHERE id = ?;
	`,s.Name, s.SNMPmib, s.SNMPaddress, s.ObjectType, s.ID).Err()
}

func (r *repository) DeleteSocket(ctx context.Context, socketID int) error {
	return r.conn.QueryRow(`
	DELETE FROM sockets WHERE id = ?;
	`, socketID).Err()
}

func (r *repository) FindSocketsByLocation(ctx context.Context, locationAddress string) ([]*manager.Socket, error) {
	rows, err := r.conn.Query(`
	SELECT id, name, mib_address, socket_type_id FROM sockets WHERE netping_address = ?;
	`, locationAddress)
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

func (r *repository) FindSocketByID(ctx context.Context, socketID int) (*manager.Socket, error) {
	socket := &manager.Socket{
		ID: socketID,
	}
	err := r.conn.QueryRow(`
	SELECT name, mib_address, netping_address, socket_type_id FROM sockets WHERE id = ?;
	`, socketID).Scan(&socket.Name, &socket.SNMPmib, &socket.SNMPaddress, &socket.ObjectType)
	return socket, err
}