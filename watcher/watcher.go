package watcher

import (
	"errors"

	"github.com/gosnmp/gosnmp"
)

type Watcher struct {
	Locations map[string]*Location
}

// Location represents a real location like 'Bishkek'.
type Location struct {
	// Address represents the netping ipv4 for the location
	Address string `json:"-"`
	// Conn is a connection for the netping at the location
	Conn *gosnmp.GoSNMP `json:"-"`
	// Sockets are "all" the Machines that are connected to the netping
	// in that location. Well not queite all of them. You have to add them
	// manualy so yeah :)
	Sockets []Socket `json:"sockets"`
}

type Socket struct {
	Name    string
	Address string
}

// These are all the types of machines that can be connector to our sockets
const (
	TypeAirConditioner = iota + 1
	TypeGenerator
	TypeHeater
)

func NewWatcher() (*Watcher, error) {
	return &Watcher{
		Locations: make(map[string]*Location),
	}, nil
}

func (w *Watcher) AddLocation(locationName, address, community string, port int) (map[string]*Location, error) {
	conn := *gosnmp.Default
	conn.Target = address
	conn.Port = uint16(port)
	conn.Community = community
	if err := conn.Connect(); err != nil {
		return nil, err
	}
	w.Locations[locationName] = &Location{
		Address: address,
		Sockets: make([]Socket, 0),
		Conn:    &conn,
	}
	return w.Locations, nil
}

func (w *Watcher) RemoveLocation(locationName string) (map[string]*Location, error) {
	delete(w.Locations, locationName)
	return w.Locations, nil
}

func (w *Watcher) AddSocket(locationName, socketName, socketMIB string) (map[string]*Location, error) {
	v, ok := w.Locations[locationName]
	if !ok {
		return nil, errors.New("watcher: no such location")
	}
	v.Sockets = append(v.Sockets, Socket{
		Name:    socketName,
		Address: socketMIB,
	})
	res, err := v.Conn.Get([]string{socketMIB})
	if err != nil {
		return nil, err
	}
	if res.Error != gosnmp.NoError {
		return nil, errors.New("watcher: something went wrong while checking socketMIB")
	}
	return w.Locations, nil
}

func (w *Watcher) RemoveSocket(locationName, socketName string) (map[string]*Location, error) {
	v, ok := w.Locations[locationName]
	if !ok {
		return nil, errors.New("watcher: there is no such location")
	}
	for i, vv := range v.Sockets {
		if vv.Name == socketName {
			v.Sockets = append(v.Sockets[:i], v.Sockets[i+1:]...)
			return w.Locations, nil
		}
	}
	return w.Locations, errors.New("watcher: there was no such socket")
}

type walkResponse struct {
	Values map[string]string `json:"values"`
}

// snmpget -v 2c -c SWITCH 192.168.0.100  .1.3.6.1.4.1.25728.8900.1.1.3.4

func (w *Watcher) Walk() (*walkResponse, error) {
	resp := &walkResponse{
		Values: make(map[string]string),
	}
	for _, v := range w.Locations {
		var oids []string
		var names []string
		for _, vv := range v.Sockets {
			oids = append(oids, vv.Address)
			names = append(names, vv.Name)
		}
		result, err := v.Conn.Get(oids)
		if err != nil {
			return nil, err
		}
		for i, vv := range result.Variables {
			if vv.Value.(int) == 0 {
				resp.Values[names[i]] = "off"
			} else {
				resp.Values[names[i]] = "on"
			}
		}
	}
	return resp, nil
}

type toggleSocketResponse struct {
	Values map[string]int `json:"values"`
}

func (w *Watcher) ToggleSocket(locationName, socketName string, value int) (*toggleSocketResponse, error) {
	resp := &toggleSocketResponse{
		Values: make(map[string]int),
	}
	l, ok := w.Locations[locationName]
	if !ok {
		return nil, errors.New("watcher: there is no such location")
	}
	s := Socket{}
	for _, v := range l.Sockets {
		if v.Name == socketName {
			s = v
			break
		}
	}
	setArgs := []gosnmp.SnmpPDU{{Name: s.Address, Value: value, Type: gosnmp.Integer}}

	result, err := l.Conn.Set(setArgs)
	if err != nil {
		return nil, err
	}
	for _, vv := range result.Variables {
		resp.Values[socketName] = vv.Value.(int)
	}
	return resp, nil
}
