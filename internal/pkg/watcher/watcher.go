package watcher

// This package is deprecated

import (
	"context"
	"errors"
	"fmt"

	"github.com/gosnmp/gosnmp"
	"github.com/rasulov-emirlan/netping-manager/internal/manager"
)

type Watcher struct {
	Locations []*Location
}

type Location struct {
	model *manager.Location
	Conn  *gosnmp.GoSNMP `json:"-"`
}

// These are all the types of machines that can be connector to our sockets
const (
	TypeAirConditioner = iota + 1
	TypeGenerator
	TypeHeater
)

func NewWatcher(l []*manager.Location) (manager.Watcher, error) {
	locations := make([]*Location, len(l))
	for i, v := range l {
		conn := *gosnmp.Default
		conn.Target = v.SNMPaddress
		conn.Port = uint16(v.SNMPport)
		conn.Community = v.SNMPcommunity
		if err := conn.Connect(); err != nil {
			return nil, err
		}
		locations[i] = &Location{
			model: v,
			Conn:  &conn,
		}
	}
	return &Watcher{
		Locations: locations,
	}, nil
}

func (w *Watcher) AddLocation(ctx context.Context, l manager.Location) ([]*manager.Location, error) {
	conn := *gosnmp.Default
	conn.Target = l.SNMPaddress
	conn.Port = uint16(l.SNMPport)
	conn.Community = l.SNMPcommunity
	if err := conn.Connect(); err != nil {
		return nil, err
	}
	w.Locations = append(w.Locations, &Location{
		model: &l,
		Conn:  &conn,
	})
	return locationsToService(w.Locations), nil
}

func (w *Watcher) RemoveLocation(ctx context.Context, locationID int) ([]*manager.Location, error) {
	check := false
	index := 0
	for i, v := range w.Locations {
		if v.model.ID == locationID {
			check = true
			index = i
			break
		}
	}
	if !check {
		return nil, errors.New("watcher: no such location")
	}
	w.Locations = append(w.Locations[:index], w.Locations[index+1:]...)
	return locationsToService(w.Locations), nil
}

func (w *Watcher) AddSocket(ctx context.Context, s manager.Socket, locationID int) ([]*manager.Location, error) {
	check := false
	index := 0
	for i, v := range w.Locations {
		if v.model.ID == locationID {
			check = true
			index = i
			break
		}
	}
	if !check {
		return nil, errors.New("watcher: no such location")
	}

	for _, v := range w.Locations[index].model.Sockets {
		if v.SNMPmib == s.SNMPmib {
			return nil, errors.New("watcher: such socket already exists")
		}
	}

	w.Locations[index].model.Sockets = append(w.Locations[index].model.Sockets, &s)
	res, err := w.Locations[index].Conn.Get([]string{s.SNMPmib})
	if err != nil {
		return nil, err
	}
	if res.Error != gosnmp.NoError {
		return nil, errors.New("watcher: something went wrong while checking socketMIB")
	}
	return locationsToService(w.Locations), nil
}

func (w *Watcher) RemoveSocket(ctx context.Context, socketID, locationID int) ([]*manager.Location, error) {
	check := false
	index := 0
	for i, v := range w.Locations {
		if v.model.ID == locationID {
			check = true
			index = i
			break
		}
	}
	if !check {
		return nil, errors.New("watcher: no such location")
	}

	check = false
	indexSocket := 0
	for i, v := range w.Locations[index].model.Sockets {
		if v.ID == socketID {
			check = true
			indexSocket = i
			break
		}
	}
	if !check {
		return nil, errors.New("watcher: no such socket")
	}

	w.Locations[index].model.Sockets = append(
		w.Locations[index].model.Sockets[:index],
		w.Locations[index].model.Sockets[indexSocket+1:]...,
	)

	return locationsToService(w.Locations), nil

}

func (w *Watcher) CheckAll(ctx context.Context) ([]*manager.Location, error) {
	for _, v := range w.Locations {
		oids := make([]string, len(v.model.Sockets))
		for i, vv := range v.model.Sockets {
			oids[i] = vv.SNMPmib
		}
		res, err := v.Conn.Get(oids)
		if err != nil {
			return nil, err
		}
		if res.Error != gosnmp.NoError {
			return nil, fmt.Errorf("watcher: gosnmp error: %d", res.Error)
		}
		if len(res.Variables) != len(v.model.Sockets) {
			return nil, errors.New("watcher: result from get is incorrect, maybe some sockets are out of work")
		}
		for i, vv := range res.Variables {
			if vv.Value.(int) == 1 {
				v.model.Sockets[i].IsON = true
				continue
			}
			v.model.Sockets[i].IsON = false
		}
	}
	return locationsToService(w.Locations), nil
}

func (w *Watcher) ToggleSocket(ctx context.Context, socketID, locationId, onOrOff int) ([]*manager.Location, error) {
	check := false
	index := 0
	for i, v := range w.Locations {
		if v.model.ID == locationId {
			check = true
			index = i
			break
		}
	}
	if !check {
		return nil, errors.New("watcher: no such location")
	}

	check = false
	indexSocket := 0
	for i, v := range w.Locations[index].model.Sockets {
		if v.ID == socketID {
			check = true
			indexSocket = i
			break
		}
	}
	if !check {
		return nil, errors.New("watcher: no such socket")
	}

	setArgs := []gosnmp.SnmpPDU{{
		Name:  w.Locations[index].model.Sockets[indexSocket].SNMPmib,
		Value: onOrOff,
		Type:  gosnmp.Integer,
	}}

	result, err := w.Locations[index].Conn.Set(setArgs)
	if err != nil {
		return nil, err
	}
	if result.Error != gosnmp.NoError {
		return nil, fmt.Errorf("watcher: gosnmp error: %d", result.Error)
	}
	if onOrOff == 1 {
		w.Locations[index].model.Sockets[indexSocket].IsON = true
	} else {
		w.Locations[index].model.Sockets[indexSocket].IsON = false
	}

	return locationsToService(w.Locations), nil
}

// snmpget -v 2c -c SWITCH 192.168.0.100  .1.3.6.1.4.1.25728.8900.1.1.3.4
