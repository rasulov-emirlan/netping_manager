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
		name, mib_address, netping_id, socket_type_id
	)
	VALUES( ?, ?, ?, ? );
`

func (r *repository) CreateSocket(ctx context.Context, s manager.Socket, locationID int) (*manager.Socket, error) {
	err := r.conn.QueryRow(createSocketSQL, s.Name, s.SNMPmib, locationID, s.ObjectType).Err()
	return &s, err
}

const updateSocketSQL = `
	UPDATE sockets SET name = ?, mib_address = ?, socket_type_id = ?
		WHERE id = ?;
`

func (r *repository) UpdateSocket(ctx context.Context, s manager.Socket) (*manager.Socket, error) {
	return &s, r.conn.QueryRow(updateSocketSQL, s.Name, s.SNMPmib, s.ObjectType, s.ID).Err()
}

const deleteSocketSQL = `
	DELETE FROM sockets WHERE id = ?;
`

func (r *repository) DeleteSocket(ctx context.Context, socketID int) error {
	return r.conn.QueryRow(deleteSocketSQL, socketID).Err()
}

const findSocketsByLocationIDSQL = `
	SELECT s.id, s.name, nl.host, s.mib_address, s.socket_type_id 
		FROM sockets AS s
	INNER JOIN netping_list AS nl
		ON s.netping_id = nl.id
	WHERE s.netping_id = ?;
`

func (r *repository) FindSocketsByLocationID(ctx context.Context, locationID int) ([]*manager.Socket, error) {
	rows, err := r.conn.Query(findSocketsByLocationIDSQL, locationID)
	if err != nil {
		return nil, err
	}

	var (
		sockets                          []*manager.Socket
		id, socketType                   int
		name, mibAddress, netpingAddress string
	)

	for rows.Next() {
		if err := rows.Scan(
			&id, &name, &netpingAddress, &mibAddress, &socketType,
		); err != nil {
			return nil, err
		}
		sockets = append(sockets, &manager.Socket{
			ID:          id,
			Name:        name,
			SNMPaddress: netpingAddress,
			SNMPmib:     mibAddress,
			ObjectType:  socketType,
		})
	}
	return sockets, nil
}

const findSocketByIDsql = `
	SELECT s.name, s.mib_address, nl.host, s.socket_type_id FROM sockets AS s
	INNER JOIN netping_list AS nl
	ON nl.id = s.netping_id WHERE s.id = ?;
`

func (r *repository) FindSocketByID(ctx context.Context, socketID int) (*manager.Socket, error) {
	socket := &manager.Socket{
		ID: socketID,
	}
	err := r.conn.QueryRow(findSocketByIDsql, socketID).Scan(&socket.Name, &socket.SNMPmib, &socket.SNMPaddress, &socket.ObjectType)
	return socket, err
}

const listAllSocketsSQL = `
	SELECT nl.id, COALESCE(s.id, 0), nl.name, COALESCE(s.name, ''), nl.host, COALESCE(s.mib_address, ''), COALESCE(s.socket_type_id, 1) 
	FROM netping_list AS nl
	LEFT JOIN sockets AS s
	ON s.netping_id = nl.id;
`

func (r *repository) ListAllSockets(ctx context.Context) ([]*manager.Location, error) {
	rows, err := r.conn.Query(listAllSocketsSQL)
	if err != nil {
		return nil, err
	}

	var (
		locationsMap = make(map[string]*manager.Location)

		netpingAddress, lname, sname, mibAddress string
		sid, lid, socketType                     int
	)

	for rows.Next() {
		if err := rows.Scan(
			&lid, &sid, &lname, &sname, &netpingAddress, &mibAddress, &socketType,
		); err != nil {
			return nil, err
		}
		if _, ok := locationsMap[lname]; !ok {
			locationsMap[lname] = &manager.Location{
				ID:            lid,
				Name:          lname,
				SNMPaddress:   netpingAddress,
				SNMPcommunity: "SWITCH",
				SNMPport:      161,
			}
		}
		if sid != 0 {
			locationsMap[lname].Sockets = append(locationsMap[lname].Sockets,
				&manager.Socket{
					ID:            sid,
					Name:          sname,
					SNMPaddress:   netpingAddress,
					SNMPcommunity: "SWITCH",
					SNMPport:      161,
					SNMPmib:       mibAddress,
					ObjectType:    socketType,
				})
		}
	}

	locations := make([]*manager.Location, len(locationsMap))
	index := 0
	for _, v := range locationsMap {
		locations[index] = v
		index++
	}
	return locations, nil
}
